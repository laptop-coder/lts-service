import { createSignal, onMount, For, Show } from "solid-js";
import { api } from "../../lib/api";
import { PERMISSIONS } from "../../lib/permissions";
import { usePermissions } from "../../lib/permissions";
import type { Subject } from "../../lib/types";

const Subjects = () => {
  const [subjects, setSubjects] = createSignal<Subject[]>([]);
  const [loading, setLoading] = createSignal(true);
  const [error, setError] = createSignal("");
  const [newSubjectName, setNewSubjectName] = createSignal("");
  const [creating, setCreating] = createSignal(false);
  const [deletingId, setDeletingId] = createSignal<number | null>(null);

  const { hasPermission } = usePermissions();

  const loadSubjects = async () => {
    try {
      const data = await api.get<{ subjects: Subject[] }>("/subjects");
      setSubjects(data.subjects);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Ошибка загрузки предметов",
      );
    } finally {
      setLoading(false);
    }
  };

  const createSubject = async (e: Event) => {
    e.preventDefault();
    if (!newSubjectName().trim()) return;

    setCreating(true);
    try {
      const formData = new URLSearchParams();
      formData.append("name", newSubjectName().trim());

      await api.post("/subjects", formData);
      setNewSubjectName("");
      await loadSubjects();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Ошибка создания предмета");
    } finally {
      setCreating(false);
    }
  };

  const deleteSubject = async (id: number) => {
    if (!confirm("Удалить предмет? Это действие нельзя отменить.")) return;

    setDeletingId(id);
    try {
      await api.delete(`/subjects/${id}`);
      await loadSubjects();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Ошибка удаления предмета");
    } finally {
      setDeletingId(null);
    }
  };

  onMount(() => {
    loadSubjects();
  });

  return (
    <div class="space-y-6">
      <div class="flex justify-between items-center">
        <h1 class="text-2xl font-bold">Управление предметами</h1>
      </div>

      <Show when={error()}>
        <div class="bg-red-100 text-red-700 p-3 rounded-lg">{error()}</div>
      </Show>

      {/* Form for subjects creating */}
      <Show when={hasPermission(PERMISSIONS.SUBJECT_CREATE)}>
        <form
          onSubmit={createSubject}
          class="bg-gray-50 p-4 rounded-lg space-y-3"
        >
          <h2 class="font-semibold">Создать новый предмет</h2>
          <div class="flex gap-2">
            <input
              type="text"
              value={newSubjectName()}
              onInput={(e) => setNewSubjectName(e.currentTarget.value)}
              placeholder="Название предмета"
              class="flex-1 px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              disabled={creating()}
            />
            <button
              type="submit"
              disabled={creating() || !newSubjectName().trim()}
              class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50"
            >
              {creating() ? "Создание..." : "Создать"}
            </button>
          </div>
        </form>
      </Show>

      {/* List of subjects */}
          <Show when={loading()}>
            <div class="text-center py-8">Загрузка...</div>
          </Show>

          <Show when={!loading() && subjects().length === 0}>
            <div class="text-center text-gray-500 py-8">
              Нет предметов. Создайте первый.
            </div>
          </Show>

          <Show when={!loading() && subjects().length > 0}>
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
                  <For each={subjects()}>
                    {(subject) => (
                      <tr class="hover:bg-gray-50">
                        <td class="px-4 py-3 text-sm text-gray-500">
                          {subject.id}
                        </td>
                        <td class="px-4 py-3 font-medium">{subject.name}</td>
                        <td class="px-4 py-3 text-right">
                          <Show
                            when={hasPermission(PERMISSIONS.SUBJECT_DELETE)}
                          >
                            <button
                              onClick={() => deleteSubject(subject.id)}
                              disabled={deletingId() === subject.id}
                              class="text-red-600 hover:text-red-800 disabled:opacity-50"
                            >
                              {deletingId() === subject.id
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
    </div>
  );
};

export default Subjects;
