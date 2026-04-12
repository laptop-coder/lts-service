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

const PostCardDetailed = (props: Props) => {
  const auth = useAuth();
  const { hasPermission } = usePermissions();
  const [loading, setLoading] = createSignal(false);
  const [error, setError] = createSignal("");

  const verifyPost = async () => {
    try {
      setLoading(true);
      await api.patch<{ posts: Post[] }>(`/posts/${props.post.id}/verify`);
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
      await api.patch<{ posts: Post[] }>(`/posts/${props.post.id}/return`);
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
        await api.delete<{}>(`/posts/${props.post.id}`);
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
      class={
        "rounded-2xl shadow-md hover:shadow-xl transition-all duration-200 overflow-hidden bg-white w-full"
      }
    >
      <div class="p-5">
        <div class="flex items-start gap-4">
          <img
            class={`w-12 h-12 rounded-full object-cover flex-shrink-0`}
            src={`/storage/storage/avatars/${props.post.author.hasAvatar ? props.post.author.id : "default"}.jpeg`}
            alt="Фото профиля"
          />

          <div class="flex-1 min-w-0">
            <div class="flex items-center justify-between flex-wrap gap-2">
              <h3 class={"text-lg font-semibold truncate text-gray-800"}>
                {props.post.name}
              </h3>
              <div class="flex items-center gap-2">
                {props.post.verified ? (
                  props.post.thingReturnedToOwner ? (
                    <span class="px-2 py-0.5 bg-green-100 text-green-700 text-xs rounded-full">
                      Найдено
                    </span>
                  ) : (
                    <span class="px-2 py-0.5 bg-red-100 text-red-700 text-xs rounded-full">
                      Не найдено
                    </span>
                  )
                ) : (
                  <span class="px-2 py-0.5 bg-yellow-100 text-yellow-700 text-xs rounded-full">
                    На модерации
                  </span>
                )}
              </div>
            </div>

            <div class="flex items-center gap-3 mt-1 text-sm text-gray-500">
              <span>
                {props.post.author.firstName} {props.post.author.lastName}
              </span>
              <span>•</span>
              <span>
                Последнее изменение: {formatDate(props.post.updatedAt)}
              </span>
            </div>

            <Show when={props.post.hasPhoto}>
              <div class="mt-3">
                <img
                  src={`/storage/storage/post_photos/${props.post.id}.jpeg`}
                  alt="Фото объявления"
                  class={"w-full h-48 object-cover rounded-xl"}
                />
              </div>
            </Show>

            <Show when={props.post.description}>
              <p class={"mt-3 text-sm line-clamp-2 text-gray-600"}>
                {props.post.description}
              </p>
            </Show>

            <div class="mt-4 flex flex-col sm:flex-row justify-between gap-3">
              <div class="flex gap-3 flex-nowrap">
                {(hasPermission(PERMISSIONS.POST_UPDATE_ANY) ||
                  (hasPermission(PERMISSIONS.POST_UPDATE_OWN) &&
                    props.post.author.id === auth.user()?.id)) &&
                  !props.post.thingReturnedToOwner && (
                    <button
                      onClick={() =>
                        (window.location.href = `/posts/${props.post.id}/edit`)
                      }
                      class="w-full sm:w-auto px-3 py-1.5 bg-blue-100 text-blue-700 text-sm rounded-lg hover:bg-blue-200 transition font-medium cursor-pointer"
                    >
                      Редактировать
                    </button>
                  )}
                {(hasPermission(PERMISSIONS.POST_DELETE_ANY) ||
                  (hasPermission(PERMISSIONS.POST_DELETE_OWN) &&
                    props.post.author.id === auth.user()?.id)) &&
                  !props.post.thingReturnedToOwner && (
                    <button
                      onClick={deletePost}
                      disabled={loading()}
                      class="w-full sm:w-auto px-3 py-1.5 bg-red-100 text-red-700 text-sm rounded-lg hover:bg-red-200 transition font-medium cursor-pointer"
                    >
                      Удалить
                    </button>
                  )}
              </div>
              <div class="flex gap-3 flex-nowrap">
                {hasPermission(PERMISSIONS.POST_VERIFY) &&
                  !props.post.verified && (
                    <button
                      onClick={verifyPost}
                      disabled={loading()}
                      class="w-full sm:w-auto px-3 py-1.5 bg-green-100 text-green-700 text-sm rounded-lg hover:bg-green-200 transition font-medium cursor-pointer"
                    >
                      Верифицировать
                    </button>
                  )}
                {(hasPermission(PERMISSIONS.POST_MARK_RETURNED_ANY) ||
                  (hasPermission(PERMISSIONS.POST_MARK_RETURNED_OWN) &&
                    props.post.author.id === auth.user()?.id)) &&
                  props.post.verified &&
                  !props.post.thingReturnedToOwner && (
                    <button
                      onClick={markReturned}
                      disabled={loading()}
                      class="w-full sm:w-auto px-3 py-1.5 bg-green-100 text-green-700 text-sm rounded-lg hover:bg-green-200 transition font-medium cursor-pointer"
                    >
                      Найдено
                    </button>
                  )}
                <button
                  onClick={() => history.back()}
                  class="w-full sm:w-auto px-3 py-1.5 bg-gray-100 text-gray-700 text-sm rounded-lg hover:bg-gray-200 transition font-medium cursor-pointer flex items-center gap-1"
                >
                  Назад
                  <svg
                    class="w-4 h-4"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M15 19l-7-7 7-7"
                    />
                  </svg>
                </button>
              </div>
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

export default PostCardDetailed;
