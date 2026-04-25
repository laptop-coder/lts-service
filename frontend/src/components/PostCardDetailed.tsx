import { Show, createSignal, onMount, onCleanup, createEffect } from "solid-js";
import type { Post } from "../lib/types";
import { usePermissions, PERMISSIONS } from "../lib/permissions";
import { api, conversationApi } from "../lib/api";
import { useAuth } from "../lib/auth";
import { formatDate } from "../lib/utils";
import { A, useNavigate } from "@solidjs/router";
import { ChevronLeft } from "lucide-solid";

interface Props {
  post: Post;
  onChange?: () => void;
}

const PostCardDetailed = (props: Props) => {
  const auth = useAuth();
  const { hasPermission } = usePermissions();
  const [loading, setLoading] = createSignal(false);
  const [error, setError] = createSignal("");
  const [contactLoading, setContactLoading] = createSignal(false);
  const [contactMessage, setContactMessage] = createSignal("");
  const [showModal, setShowModal] = createSignal(false);
  const navigate = useNavigate();

  const openModal = async () => {
    setShowModal(true);
    focusMessageInput();
  };

  const closeModal = () => {
    setShowModal(false);
    setContactMessage("");
    setError("");
    setContactLoading(false);
  };

  const handleKeyDown = (e: KeyboardEvent) => {
    if (e.key === "Escape" && showModal()) {
      closeModal();
    }
  };

  onMount(() => {
    window.addEventListener("keydown", handleKeyDown);
    onCleanup(() => {
      window.removeEventListener("keydown", handleKeyDown);
    });
  });

  const contactAuthor = async () => {
    try {
      if (!contactMessage().trim()) {
        setError("Введите сообщение");
        return;
      }
      setContactLoading(true);
      const data = await conversationApi.create(
        props.post.id,
        contactMessage(),
      );
      navigate(`/conversations/${data.conversationId}`);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Не удалось начать переписку",
      );
    } finally {
      setContactLoading(false);
    }
  };

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
      setError(
        err instanceof Error ? err.message : "Не удалось закрыть объявление",
      );
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

  let messageInputRef: HTMLInputElement | undefined;

  const focusMessageInput = () => {
    if (messageInputRef) {
      messageInputRef.focus();
    }
  };

  createEffect(() => {
    focusMessageInput();
  });

  return (
    <div
      class={
        "rounded-2xl shadow-md hover:shadow-xl transition-all duration-200 overflow-hidden bg-white w-full"
      }
    >
      <div class="p-5">
        <div class="flex items-start gap-4">
          <A
            href={`/users/${props.post.author.id}`}
            class="w-10 h-10 flex bg-gray-100 rounded-full hover:bg-gray-200 transition"
          >
            <img
              class="w-10 h-10 rounded-full object-cover border-2 border-blue-100 hover:brightness-95 transition"
              src={`/storage/storage/avatars/${props.post.author.hasAvatar ? props.post.author.id : "default"}.jpeg`}
              alt="Фото профиля"
            />
          </A>
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
              <div class="mt-7 mb-5 flex justify-center">
                <img
                  src={`/storage/storage/post_photos/${props.post.id}.jpeg`}
                  alt="Фото объявления"
                  class={"max-h-100 object-contain rounded-xl"}
                />
              </div>
            </Show>

            <Show when={props.post.description}>
              <p class={"mt-2 text-sm line-clamp-2 text-gray-600"}>
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
                      type="button"
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
                      type="button"
                      class="w-full sm:w-auto px-3 py-1.5 bg-red-100 text-red-700 text-sm rounded-lg hover:bg-red-200 transition font-medium cursor-pointer"
                    >
                      Удалить
                    </button>
                  )}
                {auth.user()?.id !== props.post.author.id && (
                  <button
                    onClick={openModal}
                    disabled={contactLoading()}
                    type="button"
                    class="w-full sm:w-auto px-3 py-1.5 bg-blue-100 text-blue-700 text-sm rounded-lg hover:bg-blue-200 transition font-medium cursor-pointer"
                  >
                    {contactLoading() ? "..." : "Связаться с автором"}
                  </button>
                )}
              </div>
              <div class="flex gap-3 flex-nowrap">
                {hasPermission(PERMISSIONS.POST_VERIFY) &&
                  !props.post.verified && (
                    <button
                      onClick={verifyPost}
                      disabled={loading()}
                      type="button"
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
                      type="button"
                      class="w-full sm:w-auto px-3 py-1.5 bg-green-100 text-green-700 text-sm rounded-lg hover:bg-green-200 transition font-medium cursor-pointer"
                    >
                      Найдено
                    </button>
                  )}
                <button
                  onClick={() => history.back()}
                  type="button"
                  class="w-full sm:w-auto px-3 py-1.5 bg-gray-100 text-gray-700 text-sm rounded-lg hover:bg-gray-200 transition font-medium cursor-pointer flex flex-row items-center"
                >
                  <ChevronLeft /> Назад
                </button>
              </div>
            </div>

            <Show when={error()}>
              <div class="mt-3 text-red-600 text-sm">{error()}</div>
            </Show>
          </div>
        </div>
      </div>

      <Show when={showModal()}>
        <div
          class="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 p-4"
          onClick={closeModal}
        >
          <div
            class="bg-white rounded-2xl shadow-2xl max-w-lg w-full max-h-[90vh] overflow-hidden"
            onClick={(e) => e.stopPropagation()}
          >
            {/* Header */}
            <div class="sticky top-0 bg-white border-b border-gray-200 px-6 py-4">
              <h2 class="text-xl font-bold text-gray-800">
                Связаться с автором
              </h2>
              <p class="text-sm text-gray-500">
                {props.post.author.firstName} {props.post.author.lastName} ·{" "}
                {props.post.name}
              </p>
            </div>

            {/* Body */}
            <div class="p-6 overflow-y-auto max-h-[calc(90vh-140px)] space-y-5 flex">
              <Show when={error()}>
                <div class="bg-red-50 border border-red-200 text-red-600 p-3 rounded-xl text-sm">
                  {error()}
                </div>
              </Show>
              <input
                ref={messageInputRef}
                disabled={contactLoading()}
                type="text"
                value={contactMessage()}
                onInput={(e) => {
                  setContactLoading(false);
                  setError("");
                  setContactMessage(e.target.value);
                }}
                onKeyDown={async (e) => {
                  if (e.key === "Enter" && !e.shiftKey) {
                    e.preventDefault();
                    if (contactLoading()) return;
                    await contactAuthor();
                  }
                }}
                placeholder="Введите сообщение..."
                class="flex-1 px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition disabled:opacity-50 disabled:cursor-not-allowed"
                required
              />
            </div>

            {/* Footer */}
            <div class="sticky bottom-0 bg-white border-t border-gray-200 px-6 py-4 flex justify-end gap-3">
              <button
                onClick={closeModal}
                class="px-4 py-2 bg-gray-100 text-gray-700 rounded-xl hover:bg-gray-200 transition font-medium cursor-pointer"
              >
                Отмена
              </button>
              <button
                onClick={contactAuthor}
                disabled={contactLoading()}
                class="px-4 py-2 bg-blue-600 text-white rounded-xl hover:bg-blue-700 transition font-medium disabled:opacity-50 cursor-pointer disabled:cursor-not-allowed"
              >
                {contactLoading() ? "Отправка..." : "Отправить"}
              </button>
            </div>
          </div>
        </div>
      </Show>
    </div>
  );
};

export default PostCardDetailed;
