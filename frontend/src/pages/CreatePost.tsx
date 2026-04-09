import { createSignal, Show } from "solid-js";
import { usePermissions, PERMISSIONS } from "../lib/permissions";
import { api } from "../lib/api";
import { useNavigate } from "@solidjs/router";
import type { Post } from "../lib/types";

const CreatePost = () => {
  const [name, setName] = createSignal("");
  const [description, setDescription] = createSignal("");
  const [photo, setPhoto] = createSignal<File | null>(null);
  const [photoPreview, setPhotoPreview] = createSignal<string | null>(null);
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
    if (photo()) formData.append("photo", photo()!);

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

  const handlePhotoChange = (e: Event) => {
    const input = e.currentTarget as HTMLInputElement;
    const file = input.files?.[0];
    if (file) {
      setPhoto(file);
      const preview = URL.createObjectURL(file);
      setPhotoPreview(preview);
    }
  };

  const removePhoto = () => {
    setPhoto(null);
    if (photoPreview()) {
      URL.revokeObjectURL(photoPreview()!);
      setPhotoPreview(null);
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
                class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition"
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
                class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition min-h-[140px] max-h-[600px]"
              />
              <p class="text-xs text-gray-500 mt-1">Чем подробнее, тем лучше</p>
            </div>

            {/* Photo upload */}
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">
                Фото
              </label>

              <Show when={!photoPreview()}>
                <label class="flex flex-col items-center justify-center w-full h-32 border-2 border-dashed border-gray-300 rounded-xl cursor-pointer hover:border-blue-500 transition">
                  <div class="flex flex-col items-center justify-center pt-5 pb-6">
                    <svg
                      class="w-8 h-8 text-gray-400 mb-2"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
                      ></path>
                    </svg>
                    <p class="text-sm text-gray-500">
                      Нажмите для загрузки фото
                    </p>
                    <p class="text-xs text-gray-400 mt-1">
                      JPEG, PNG, WebP, GIF (макс. 10MB)
                    </p>
                  </div>
                  <input
                    type="file"
                    accept="image/jpeg,image/png,image/webp,image/gif"
                    onChange={handlePhotoChange}
                    class="hidden"
                  />
                </label>
              </Show>

              <Show when={photoPreview()}>
                <div class="relative">
                  <img
                    src={photoPreview()!}
                    alt="Preview"
                    class="w-full h-48 object-cover rounded-xl"
                  />
                  <button
                    type="button"
                    onClick={removePhoto}
                    class="absolute top-2 right-2 p-1 bg-red-500 text-white rounded-full hover:bg-red-600 transition"
                  >
                    <svg
                      class="w-5 h-5"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M6 18L18 6M6 6l12 12"
                      ></path>
                    </svg>
                  </button>
                </div>
              </Show>
              <p class="text-xs text-gray-500 mt-1">
                Вы можете добавить одно фото
              </p>
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
                class="flex-1 px-4 py-2 bg-blue-600 text-white rounded-xl hover:bg-blue-700 transition font-medium disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
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
