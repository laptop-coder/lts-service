import { createSignal, onMount, For, Show } from "solid-js";
import { api } from "../../lib/api";
import { PERMISSIONS } from "../../lib/permissions";
import { usePermissions } from "../../lib/permissions";
import type { StudentGroup } from "../../lib/types";

const StudentGroupsManagement = () => {
  const [studentGroups, setStudentGroups] = createSignal<StudentGroup[]>([]);
  const [loading, setLoading] = createSignal(true);
  const [error, setError] = createSignal("");
  const [newStudentGroupName, setNewStudentGroupName] = createSignal("");
  const [creating, setCreating] = createSignal(false);
  const [deletingId, setDeletingId] = createSignal<number | null>(null);

  const { hasPermission } = usePermissions();

  const loadStudentGroups = async () => {
    try {
      const data = await api.get<{ studentGroups: StudentGroup[] }>(
        "/student_groups",
      );
      setStudentGroups(data.studentGroups);
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
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : "Ошибка удаления учебной группы (класса)",
      );
    } finally {
      setDeletingId(null);
    }
  };

  onMount(() => {
    loadStudentGroups();
  });

  return (
    <div class="space-y-6">
      <div class="flex justify-between items-center">
        <h1 class="text-2xl font-bold">Управление учебными группами</h1>
      </div>

      <Show when={error()}>
        <div class="bg-red-100 text-red-700 p-3 rounded-lg">{error()}</div>
      </Show>

      {/* Form for student group creating */}
      <Show when={hasPermission(PERMISSIONS.STUDENT_GROUP_CREATE)}>
        <form
          onSubmit={createStudentGroup}
          class="bg-gray-50 p-4 rounded-lg space-y-3"
        >
          <h2 class="font-semibold">Создать новую группу (класс)</h2>
          <div class="flex gap-2">
            <input
              type="text"
              value={newStudentGroupName()}
              onInput={(e) => setNewStudentGroupName(e.currentTarget.value)}
              placeholder="Название группы"
              class="flex-1 px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              disabled={creating()}
            />
            <button
              type="submit"
              disabled={creating() || !newStudentGroupName().trim()}
              class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50"
            >
              {creating() ? "Создание..." : "Создать"}
            </button>
          </div>
        </form>
      </Show>

      {/* List of student groups */}
      {hasPermission(PERMISSIONS.STUDENT_GROUP_READ_ANY) && (
        <>
          <Show when={loading()}>
            <div class="text-center py-8">Загрузка...</div>
          </Show>

          <Show when={!loading() && studentGroups().length === 0}>
            <div class="text-center text-gray-500 py-8">
              Нет групп. Создайте первую.
            </div>
          </Show>

          <Show when={!loading() && studentGroups().length > 0}>
            <div class="bg-white rounded-lg shadow overflow-hidden">
              <table class="w-full">
                <thead class="bg-gray-50">
                  <tr>
                    <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">
                      ID
                    </th>
                    <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">
                      Название
                    </th>
                    <th class="px-4 py-3 text-right text-sm font-medium text-gray-500">
                      Действия
                    </th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-gray-200">
                  <For each={studentGroups()}>
                    {(studentGroup) => (
                      <tr class="hover:bg-gray-50">
                        <td class="px-4 py-3 text-sm text-gray-500">
                          {studentGroup.id}
                        </td>
                        <td class="px-4 py-3 font-medium">
                          {studentGroup.name}
                        </td>
                        <td class="px-4 py-3 text-right">
                          <Show
                            when={hasPermission(
                              PERMISSIONS.STUDENT_GROUP_DELETE,
                            )}
                          >
                            <button
                              onClick={() =>
                                deleteStudentGroup(studentGroup.id)
                              }
                              disabled={deletingId() === studentGroup.id}
                              class="text-red-600 hover:text-red-800 disabled:opacity-50"
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
          </Show>
        </>
      )}
    </div>
  );
};

export default StudentGroupsManagement;
