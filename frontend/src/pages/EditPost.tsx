import { createSignal, Show, onMount } from "solid-js";
import { usePermissions, PERMISSIONS } from "../lib/permissions";
import { api } from "../lib/api";
import { useAuth } from "../lib/auth";
import { useNavigate, useParams } from "@solidjs/router";
import type { Post } from "../lib/types";
import { X, Image } from "lucide-solid";

const EditPost = () => {
  const params = useParams();
  const navigate = useNavigate();
  const [postAuthorId, setPostAuthorId] = createSignal<string | null>(null);
  const [name, setName] = createSignal("");
  const [description, setDescription] = createSignal("");
  const [error, setError] = createSignal("");
  const [initialLoading, setInitialLoading] = createSignal(true);
  const [loading, setLoading] = createSignal(false);
  const { hasPermission } = usePermissions();
  const auth = useAuth();
  const [hasPhoto, setHasPhoto] = createSignal(false);
  const [photoPreview, setPhotoPreview] = createSignal<string | null>(null);
  const [newPhoto, setNewPhoto] = createSignal<File | null>(null);

  const loadPost = async () => {
    try {
      const data = await api.get<{ post: Post }>(`/posts/${params.id}`);
      setPostAuthorId(data.post.author.id);
      setName(data.post.name);
      setDescription(data.post.description || "");
      setHasPhoto(data.post.hasPhoto);
      if (data.post.hasPhoto) {
        setPhotoPreview(`/storage/storage/post_photos/${data.post.id}.jpeg`);
      }
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Не удалось загрузить объявление",
      );
    } finally {
      setInitialLoading(false);
    }
  };

  const handlePhotoChange = (e: Event) => {
    const input = e.currentTarget as HTMLInputElement;
    const file = input.files?.[0];
    if (file) {
      setNewPhoto(file);
      const preview = URL.createObjectURL(file);
      setPhotoPreview(preview);
    }
  };

  const removePhoto = async () => {
    if (hasPhoto() && !newPhoto()) {
      try {
        await api.delete(`/posts/${params.id}/photo`);
        setHasPhoto(false);
        setPhotoPreview(null);
      } catch (err) {
        setError(
          err instanceof Error ? err.message : "Не удалось удалить фото",
        );
      }
    } else {
      setNewPhoto(null);
      setPhotoPreview(
        hasPhoto() ? `/storage/storage/post_photos/${params.id}.jpeg` : null,
      );
    }
  };

  onMount(() => {
    loadPost();
  });

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    const formData = new URLSearchParams();
    formData.append("name", name());
    if (description()?.trim()) formData.append("description", description());

    try {
      await api.patch(`/posts/${params.id}`, formData);
      if (newPhoto()) {
        const photoFormData = new FormData();
        photoFormData.append("photo", newPhoto()!);
        await api.put(`/posts/${params.id}/photo`, photoFormData);
      }
      navigate("/");
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Не удалось обновить объявление",
      );
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <Show when={initialLoading()}>
        <div class="text-center py-8">Загрузка...</div>;
      </Show>
      <Show when={!initialLoading()}>
        {(hasPermission(PERMISSIONS.POST_UPDATE_ANY) ||
          (hasPermission(PERMISSIONS.POST_UPDATE_OWN) &&
            postAuthorId() &&
            postAuthorId() === auth.user()?.id)) && (
          <div class="max-w-2xl mx-auto">
            <h1 class="text-2xl font-bold text-gray-800 text-center mb-6">
              Редактировать объявление
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
                  class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition"
                  required
                />
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">
                  Описание
                </label>
                <textarea
                  value={description()}
                  onInput={(e) => setDescription(e.currentTarget.value)}
                  rows={5}
                  class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition min-h-[140px] max-h-[600px]"
                />
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">
                  Фото
                </label>

                <Show when={!photoPreview()}>
                  {(hasPermission(PERMISSIONS.POST_PHOTO_UPDATE_ANY) ||
                    (hasPermission(PERMISSIONS.POST_PHOTO_UPDATE_OWN) &&
                      postAuthorId() === auth.user()?.id)) && (
                    <label class="flex flex-col items-center justify-center w-full h-32 border-2 border-dashed border-gray-300 rounded-xl cursor-pointer hover:border-blue-500 transition">
                      <div class="flex flex-col items-center justify-center pt-5 pb-6">
                        <Image />
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
                  )}
                </Show>

                <Show when={photoPreview()}>
                  <div class="relative">
                    <img
                      src={photoPreview()!}
                      alt="Preview"
                      class="w-full h-48 object-cover rounded-xl"
                    />
                    {(hasPermission(PERMISSIONS.POST_PHOTO_DELETE_ANY) ||
                      (hasPermission(PERMISSIONS.POST_PHOTO_DELETE_OWN) &&
                        postAuthorId() === auth.user()?.id)) && (
                      <button
                        type="button"
                        onClick={removePhoto}
                        class="absolute top-2 right-2 p-1 bg-red-500 text-white rounded-full hover:bg-red-600 transition cursor-pointer disabled:cursor-not-allowed"
                      >
                        <X />
                      </button>
                    )}
                  </div>
                </Show>
                <p class="text-xs text-gray-500 mt-1">
                  Вы можете добавить или заменить фото
                </p>
              </div>

              <Show when={error()}>
                <div class="bg-red-50 text-red-600 p-3 rounded-xl text-sm border border-red-200">
                  {error()}
                </div>
              </Show>

              <div class="flex gap-3 pt-2">
                <button
                  type="button"
                  onClick={() => navigate("/")}
                  class="flex-1 px-4 h-40 bg-gray-100 text-gray-700 rounded-xl hover:bg-gray-200 transition font-medium cursor-pointer"
                >
                  Отмена
                </button>
                <button
                  type="submit"
                  disabled={loading()}
                  class="flex-1 px-4 h-40 bg-blue-600 text-white rounded-xl hover:bg-blue-700 transition font-medium disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
                >
                  Сохранить
                </button>
              </div>
            </form>
          </div>
        )}
      </Show>
    </>
  );
};

export default EditPost;
