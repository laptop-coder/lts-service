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
    <div class="space-y-6 p-4">
      <div class="mb-6">
        <h1 class="text-3xl font-bold text-gray-800">Управление предметами</h1>
        <p class="text-gray-500 mt-1">Список учебных дисциплин</p>
      </div>

      <Show when={error()}>
        <div class="bg-red-50 border border-red-200 text-red-600 p-3 rounded-xl">
          {error()}
        </div>
      </Show>

      {/* Form for subjects creating */}
      <Show when={hasPermission(PERMISSIONS.SUBJECT_CREATE)}>
        <div class="bg-white rounded-2xl shadow-lg p-6">
          <h2 class="text-lg font-semibold text-gray-800 mb-4">
            Добавить предмет
          </h2>
          <form onSubmit={createSubject} class="flex gap-3">
            <input
              type="text"
              value={newSubjectName()}
              onInput={(e) => setNewSubjectName(e.currentTarget.value)}
              placeholder="Название предмета"
              class="flex-1 px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition disabled:opacity-50"
              disabled={creating()}
            />
            <button
              type="submit"
              disabled={creating() || !newSubjectName().trim()}
              class="px-5 py-2 bg-blue-600 text-white rounded-xl hover:bg-blue-700 disabled:opacity-50 transition cursor-pointer disabled:cursor-not-allowed font-medium"
            >
              {creating() ? "Создание..." : "Создать"}
            </button>
          </form>
        </div>
      </Show>

      {/* List of subjects */}
      <Show when={loading()}>
        <div class="flex justify-center py-16">
          <div class="text-gray-500">Загрузка...</div>
        </div>
      </Show>

      <Show when={!loading() && subjects().length === 0}>
        <div class="text-center py-16">
          <div class="text-5xl mb-3">📚</div>
          <p class="text-gray-500">Нет предметов</p>
          <p class="text-gray-400 text-sm mt-1">Создайте первый</p>
        </div>
      </Show>

      <Show when={!loading() && subjects().length > 0}>
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
                <For each={subjects()}>
                  {(subject) => (
                    <tr class="hover:bg-gray-50 transition">
                      <td class="px-6 py-4 text-sm text-gray-500 font-mono">
                        {subject.id}
                      </td>
                      <td class="px-6 py-4 font-medium text-gray-800">
                        {subject.name}
                      </td>
                      <td class="px-6 py-4 text-right">
                        <Show when={hasPermission(PERMISSIONS.SUBJECT_DELETE)}>
                          <button
                            onClick={() => deleteSubject(subject.id)}
                            disabled={deletingId() === subject.id}
                            class="text-red-600 hover:text-red-800 disabled:opacity-50 transition cursor-pointer disabled:cursor-not-allowed font-medium"
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
        </div>
      </Show>
    </div>
  );
};

export default Subjects;
