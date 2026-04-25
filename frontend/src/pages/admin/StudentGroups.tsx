import { createSignal, createEffect, onMount, For, Show } from "solid-js";
import { api } from "../../lib/api";
import { PERMISSIONS } from "../../lib/permissions";
import { usePermissions } from "../../lib/permissions";
import type { StudentGroup } from "../../lib/types";
import Pagination from '../../components/Pagination'

const StudentGroups = () => {
  const [studentGroups, setStudentGroups] = createSignal<StudentGroup[]>([]);
  const [loading, setLoading] = createSignal(true);
  const [error, setError] = createSignal("");
  const [newStudentGroupName, setNewStudentGroupName] = createSignal("");
  const [creating, setCreating] = createSignal(false);
  const [deletingId, setDeletingId] = createSignal<number | null>(null);
  const [page, setPage] = createSignal(0);
  const [hasMore, setHasMore] = createSignal(true);

  const { hasPermission } = usePermissions();

  let inputRef: HTMLInputElement | undefined;
  const focusInput = () => {
    if (inputRef) {
      inputRef.focus();
    }
  };

  createEffect(() => {
    focusInput();
  });

  const limit = 30

  createEffect(() => {
    page()
    loadStudentGroups()
  })

  const loadStudentGroups = async () => {
    try {
      const data = await api.get<{ studentGroups: StudentGroup[] }>(
`/student_groups?limit=${limit+1}&offset=${page() * limit}`
      );
      setHasMore(data.studentGroups.length > limit);
      setStudentGroups(data.studentGroups.slice(0, limit));
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : "Ошибка загрузки учебных групп (классов)",
      );
    } finally {
      setLoading(false);
    }
  };

  const createStudentGroup = async (e: Event) => {
    e.preventDefault();
    if (!newStudentGroupName().trim()) return;

    setCreating(true);
    try {
      const formData = new URLSearchParams();
      formData.append("name", newStudentGroupName().trim());

      await api.post("/student_groups", formData);
      setNewStudentGroupName("");
      await loadStudentGroups();
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : "Ошибка создания учебной группы (класса)",
      );
    } finally {
      setCreating(false);
      focusInput();
    }
  };

  const deleteStudentGroup = async (id: number) => {
    if (
      !confirm("Удалить учебную группу (класс)? Это действие нельзя отменить.")
    )
      return;

    setDeletingId(id);
    try {
      await api.delete(`/student_groups/${id}`);
      await loadStudentGroups();
      if (studentGroups().length === 0 && page() > 0) {
        setPage(prev => prev - 1)
      }
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : "Ошибка удаления учебной группы (класса)",
      );
    } finally {
      setDeletingId(null);
      focusInput();
    }
  };

  onMount(() => {
    loadStudentGroups();
  });

  return (
    <div class="space-y-6 p-4">
      <div class="mb-6">
        <h1 class="text-3xl font-bold text-gray-800">
          Управление учебными группами
        </h1>
        <p class="text-gray-500 mt-1">Список классов и учебных групп</p>
      </div>

      <Show when={error()}>
        <div class="bg-red-50 border border-red-200 text-red-600 p-3 rounded-xl">
          {error()}
        </div>
      </Show>

      {/* Form for student group creating */}
      <Show when={hasPermission(PERMISSIONS.STUDENT_GROUP_CREATE)}>
        <div class="bg-white rounded-2xl shadow-lg p-6">
          <h2 class="text-lg font-semibold text-gray-800 mb-4">
            Добавить группу
          </h2>
          <form onSubmit={createStudentGroup} class="flex gap-3">
            <input
              ref={inputRef}
              type="text"
              value={newStudentGroupName()}
              onInput={(e) => setNewStudentGroupName(e.currentTarget.value)}
              placeholder="Название группы (например, 11А)"
              class="flex-1 px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition disabled:opacity-50"
              disabled={creating()}
            />
            <button
              type="submit"
              disabled={creating() || !newStudentGroupName().trim()}
              class="px-5 py-2 bg-blue-600 text-white rounded-xl hover:bg-blue-700 disabled:opacity-50 transition cursor-pointer disabled:cursor-not-allowed font-medium"
            >
              {creating() ? "Создание..." : "Создать"}
            </button>
          </form>
        </div>
      </Show>

      {/* List of student groups */}
      <Show when={loading()}>
        <div class="flex justify-center py-16">
          <div class="text-gray-500">Загрузка...</div>
        </div>
      </Show>

      <Show when={!loading() && studentGroups().length === 0}>
        <div class="text-center py-16">
          <div class="text-5xl mb-3">🎓</div>
          <p class="text-gray-500">Нет учебных групп</p>
          <p class="text-gray-400 text-sm mt-1">Создайте первую</p>
        </div>
      </Show>

      <Show when={!loading() && studentGroups().length > 0}>
        <div class="bg-white rounded-2xl shadow-lg overflow-hidden">
          <div class="overflow-x-auto">
            <table class="w-full">
              <thead class="bg-gray-50 border-b border-gray-200">
                <tr>
                  <th class="px-6 py-4 text-left text-sm font-semibold text-gray-600">
                    ID
                  </th>
                  <th class="px-6 py-4 text-left text-sm font-semibold text-gray-600">
                    Название
                  </th>
                  <th class="px-6 py-4 text-right text-sm font-semibold text-gray-600">
                    Действия
                  </th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-100">
                <For each={studentGroups()}>
                  {(studentGroup) => (
                    <tr class="hover:bg-gray-50 transition">
                      <td class="px-6 py-4 text-sm text-gray-500 font-mono">
                        {studentGroup.id}
                      </td>
                      <td class="px-6 py-4 font-medium text-gray-800">
                        {studentGroup.name}
                      </td>
                      <td class="px-6 py-4 text-right">
                        <Show
                          when={hasPermission(PERMISSIONS.STUDENT_GROUP_DELETE)}
                        >
                          <button
                            onClick={() => deleteStudentGroup(studentGroup.id)}
                            disabled={deletingId() === studentGroup.id}
                            class="text-red-600 hover:text-red-800 disabled:opacity-50 transition cursor-pointer disabled:cursor-not-allowed font-medium"
                          >
                            {deletingId() === studentGroup.id
                              ? "Удаление..."
                              : "Удалить"}
                          </button>
                        </Show>
                      </td>
                    </tr>
                  )}
                </For>
              </tbody>
            </table>
          </div>
        </div>
      </Show>
      <Pagination
        page={page()}
        hasMore={hasMore()}
        onPrev={() => setPage((prev) => prev - 1)}
        onNext={() => setPage((prev) => prev + 1)}
      />
    </div>
  );
};

export default StudentGroups;
