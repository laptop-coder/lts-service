import { createSignal } from "solid-js";
import { usePermissions, PERMISSIONS } from "../lib/permissions";
import { api } from "../lib/api";
import { useNavigate } from "@solidjs/router";
import type { Post } from "../lib/types";

// TODO: add photo support
const CreatePost = () => {
  const [name, setName] = createSignal("");
  const [description, setDescription] = createSignal("");
  const [error, setError] = createSignal("");
  const [loading, setLoading] = createSignal(false);
  const { hasPermission } = usePermissions();
  const navigate = useNavigate();

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    setError("");
    setLoading(true);
    const formData = new FormData();
    formData.append("name", name());
    if (description()?.trim()) formData.append("description", description());

    try {
      await api.post<{ posts: Post }>("/posts", formData);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Не удалось создать объявление",
      );
    } finally {
      setLoading(false);
      navigate("/");
    }
  };

  return (
    <>
      {hasPermission(PERMISSIONS.POST_CREATE) && (
        <div class="max-w-2xl mx-auto">
          <h1 class="text-2xl font-bold text-gray-800 text-center mb-6">
            Создать объявление
          </h1>

          <form
            onSubmit={handleSubmit}
            class="bg-white rounded-2xl shadow-lg p-6 space-y-5"
          >
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">
                Название *
              </label>
              <input
                type="text"
                value={name()}
                onInput={(e) => setName(e.currentTarget.value)}
                placeholder="Например: синяя шапка, чёрный рюкзак"
                class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-200 outline-none transition"
                required
              />
              <p class="text-xs text-gray-500 mt-1">
                Коротко опишите, что потеряли
              </p>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">
                Описание
              </label>
              <textarea
                value={description()}
                onInput={(e) => setDescription(e.currentTarget.value)}
                placeholder="Где и когда потеряли, особые приметы..."
                rows={5}
                class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-200 outline-none transition resize-none"
              />
              <p class="text-xs text-gray-500 mt-1">Чем подробнее, тем лучше</p>
            </div>

            {error() && (
              <div class="bg-red-50 text-red-600 p-3 rounded-xl text-sm border border-red-200">
                {error()}
              </div>
            )}

            <div class="flex gap-3 pt-2">
              <button
                type="button"
                onClick={() => navigate("/")}
                class="flex-1 px-4 py-2 bg-gray-100 text-gray-700 rounded-xl hover:bg-gray-200 transition font-medium cursor-pointer"
              >
                Отмена
              </button>
              <button
                type="submit"
                disabled={loading()}
                class="flex-1 px-4 py-2 bg-blue-50 text-blue-700 rounded-xl hover:bg-blue-100 transition font-medium disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
              >
                {loading() ? "Отправка..." : "Отправить"}
              </button>
            </div>
          </form>
        </div>
      )}
    </>
  );
};

export default CreatePost;
