import { Show, createSignal } from "solid-js";
import type { Post } from "../lib/types";
import { usePermissions, PERMISSIONS } from "../lib/permissions";
import { api } from "../lib/api";
import { useAuth } from "../lib/auth";
import { formatDate } from "../lib/utils";

interface Props {
  post: Post;
  onChange?: () => void;
}

const PostCardCompact = (props: Props) => {
  const auth = useAuth();
  const { post } = props;
  const { hasPermission } = usePermissions();
  const [loading, setLoading] = createSignal(false);
  const [error, setError] = createSignal("");

  const verifyPost = async () => {
    try {
      setLoading(true);
      await api.patch<{ posts: Post[] }>(`/posts/${post.id}/verify`);
      props.onChange?.();
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : "Не удалось верифицировать объявление",
      );
    } finally {
      setLoading(false);
    }
  };

  const markReturned = async () => {
    try {
      setLoading(true);
      await api.patch<{ posts: Post[] }>(`/posts/${post.id}/return`);
      props.onChange?.();
    } catch (err) {
      setError("Не удалось закрыть объявление");
    } finally {
      setLoading(false);
    }
  };

  const deletePost = async () => {
    if (confirm("Удалить объявление? Это действие необратимо.")) {
      try {
        setLoading(true);
        await api.delete<{}>(`/posts/${post.id}`);
        props.onChange?.();
      } catch (err) {
        setError(
          err instanceof Error ? err.message : "Не удалось удалить объявление",
        );
      } finally {
        setLoading(false);
      }
    }
  };

  return (
    <div
      class={`rounded-2xl shadow-md hover:shadow-xl transition-all duration-200 overflow-hidden ${post.thingReturnedToOwner ? "bg-gray-100 opacity-75" : "bg-white"}`}
    >
      <div class="p-5">
        <div class="flex items-start gap-4">
          <img
            class={`w-12 h-12 rounded-full object-cover flex-shrink-0 ${post.thingReturnedToOwner ? "grayscale" : ""}`}
            src={`/storage/storage/avatars/${post.author.hasAvatar ? post.author.id : "default"}.jpeg`}
            alt="Фото профиля"
          />

          <div class="flex-1 min-w-0">
            <div class="flex items-center justify-between flex-wrap gap-2">
              <h3
                class={`text-lg font-semibold truncate ${post.thingReturnedToOwner ? "text-gray-500 line-through" : "text-gray-800"}`}
              >
                {post.name}
              </h3>
              <div class="flex items-center gap-2">
                {post.thingReturnedToOwner && (
                  <span class="px-2 py-0.5 bg-green-100 text-green-700 text-xs rounded-full">
                    Найдено
                  </span>
                )}
              </div>
            </div>

            <div class="flex items-center gap-3 mt-1 text-sm text-gray-500">
              <span>
                {post.author.firstName} {post.author.lastName}
              </span>
              <span>•</span>
              <span>Последнее изменение: {formatDate(post.updatedAt)}</span>
            </div>

            <Show when={post.hasPhoto}>
              <div class="mt-3">
                <img
                  src={`/storage/storage/post_photos/${post.id}.jpeg`}
                  alt="Фото объявления"
                  class={`w-full h-48 object-cover rounded-xl ${post.thingReturnedToOwner ? "grayscale opacity-50" : ""}`}
                />
              </div>
            </Show>

            <Show when={post.description}>
              <p
                class={`mt-3 text-sm line-clamp-2 ${post.thingReturnedToOwner ? "text-gray-400" : "text-gray-600"}`}
              >
                {post.description}
              </p>
            </Show>

            <div class="flex items-center gap-3 mt-4">
              {(hasPermission(PERMISSIONS.POST_UPDATE_ANY) ||
                (hasPermission(PERMISSIONS.POST_UPDATE_OWN) &&
                  post.author.id === auth.user()?.id)) &&
                !post.thingReturnedToOwner && (
                  <a
                    href={`/posts/${post.id}/edit`}
                    class="px-3 py-1.5 bg-blue-100 text-blue-700 text-sm rounded-lg hover:bg-blue-200 transition font-medium cursor-pointer"
                  >
                    Редактировать
                  </a>
                )}
              {hasPermission(PERMISSIONS.POST_VERIFY) && !post.verified && (
                <button
                  onClick={verifyPost}
                  disabled={loading()}
                  class="px-3 py-1.5 bg-green-100 text-green-700 text-sm rounded-lg hover:bg-green-200 transition font-medium cursor-pointer"
                >
                  Верифицировать
                </button>
              )}
              {(hasPermission(PERMISSIONS.POST_MARK_RETURNED_ANY) ||
                (hasPermission(PERMISSIONS.POST_MARK_RETURNED_OWN) &&
                  post.author.id === auth.user()?.id)) &&
                post.verified &&
                !post.thingReturnedToOwner && (
                  <button
                    onClick={markReturned}
                    disabled={loading()}
                    class="px-3 py-1.5 bg-green-100 text-green-700 text-sm rounded-lg hover:bg-green-200 transition font-medium cursor-pointer"
                  >
                    Найдено
                  </button>
                )}

              {(hasPermission(PERMISSIONS.POST_DELETE_ANY) ||
                (hasPermission(PERMISSIONS.POST_DELETE_OWN) &&
                  post.author.id === auth.user()?.id)) && !post.thingReturnedToOwner && (
                <button
                  onClick={deletePost}
                  disabled={loading()}
                  class="px-3 py-1.5 bg-red-100 text-red-700 text-sm rounded-lg hover:bg-red-200 transition font-medium cursor-pointer"
                >
                  Удалить
                </button>
              )}
            </div>

            <Show when={error()}>
              <div class="mt-3 text-red-600 text-sm">{error()}</div>
            </Show>
          </div>
        </div>
      </div>
    </div>
  );
};

export default PostCardCompact;
