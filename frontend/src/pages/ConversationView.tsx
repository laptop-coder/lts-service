import { createSignal, onMount, For, Show, createEffect } from "solid-js";
import { A } from "@solidjs/router";
import { useParams, useNavigate } from "@solidjs/router";
import { conversationApi } from "../lib/api";
import { Conversation, Message } from "../lib/types";
import { useAuth } from "../lib/auth";
import { refreshUnreadMessagesCount } from "../lib/store";
import { ChevronLeft } from "lucide-solid";

const ConversationView = () => {
  const params = useParams();
  const navigate = useNavigate();
  const auth = useAuth();

  const [conversation, setConversation] = createSignal<Conversation | null>(
    null,
  );
  const [messages, setMessages] = createSignal<Message[]>([]);
  const [newMessage, setNewMessage] = createSignal("");
  const [loading, setLoading] = createSignal(true);
  const [sending, setSending] = createSignal(false);
  const [error, setError] = createSignal("");

  let messagesEndRef: HTMLDivElement | undefined;
  let messageInputRef: HTMLInputElement | undefined;

  const loadConversation = async () => {
    try {
      const data = await conversationApi.getById(params.id!);
      setConversation(data.conversation);
      setMessages(data.conversation.messages);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Ошибка загрузки переписки",
      );
    } finally {
      setLoading(false);
    }
  };

  const scrollToBottom = () => {
    messagesEndRef?.scrollIntoView({ behavior: "smooth" });
  };

  onMount(async () => {
    if (!params.id) return;
    await loadConversation();
    scrollToBottom();
    await conversationApi.markAsRead(params.id);
    await refreshUnreadMessagesCount();
  });

  const focusMessageInput = () => {
    if (messageInputRef) {
      messageInputRef.focus();
    }
  };

  createEffect(() => {
    focusMessageInput();
  });

  const sendMessage = async (e: Event) => {
    e.preventDefault();
    if (!newMessage().trim()) return;

    setSending(true);
    try {
      await conversationApi.sendMessage(params.id!, newMessage().trim());
      setNewMessage("");
      await loadConversation();
      scrollToBottom();
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Ошибка отправки сообщения",
      );
    } finally {
      setSending(false);
      focusMessageInput();
    }
  };

  const otherUser = () => conversation()?.otherUser;
  const post = () => conversation()?.post;

  return (
    <div class="max-w-4xl mx-auto h-[calc(100vh-120px)] flex flex-col bg-white rounded-2xl shadow-lg overflow-hidden">
      {/* Header */}
      <div class="border-b border-gray-200 p-4 flex items-center gap-3">
        <button
          onClick={() => navigate("/conversations")}
          type="button"
          class="text-gray-500 hover:text-gray-700 cursor-pointer flex flex-row"
        >
          <ChevronLeft /> Назад
        </button>
        <Show when={otherUser()}>
          <A
            href={`/users/${otherUser()!.id}`}
            class="flex flex-1 transition gap-2 rounded-2xl"
          >
            <img
              src={`/storage/storage/avatars/${otherUser()!.hasAvatar ? otherUser()!.id : "default"}.jpeg`}
              alt={`Фото профиля пользователя ${otherUser()!.firstName} ${otherUser()!.lastName}`}
              class="w-10 h-10 rounded-full object-cover"
            />

            <div>
              <h2 class="font-semibold text-gray-800">
                {otherUser()!.firstName} {otherUser()!.lastName}
              </h2>
              <p class="text-sm text-gray-500">{post()?.name}</p>
            </div>
          </A>
        </Show>
      </div>

      {/* Messages */}
      <div class="flex-1 overflow-y-auto p-4 space-y-3">
        <Show when={loading()}>
          <div class="text-center py-8 text-gray-500">Загрузка...</div>
        </Show>

        <Show when={error()}>
          <div class="bg-red-50 text-red-600 p-3 rounded-xl">{error()}</div>
        </Show>

        <For each={messages()}>
          {(msg) => {
            const isOwn = msg.senderId === auth.user()?.id;
            return (
              <div class={`flex ${isOwn ? "justify-end" : "justify-start"}`}>
                <div class={`max-w-[70%] ${isOwn ? "order-2" : ""}`}>
                  <div
                    class={`rounded-2xl px-4 py-2 ${
                      isOwn
                        ? "bg-blue-600 text-white"
                        : "bg-gray-100 text-gray-800"
                    }`}
                  >
                    <p class="text-sm">{msg.content}</p>
                  </div>
                  <p class="text-xs text-gray-400 mt-1">
                    {new Date(msg.createdAt).toLocaleTimeString("ru", {
                      hour: "2-digit",
                      minute: "2-digit",
                    })}
                  </p>
                </div>
              </div>
            );
          }}
        </For>
        <div ref={messagesEndRef} />
      </div>

      {/* Input */}
      <form
        onSubmit={sendMessage}
        class="border-t border-gray-200 p-4 flex gap-2"
      >
        <input
          ref={messageInputRef}
          type="text"
          value={newMessage()}
          onInput={(e) => setNewMessage(e.currentTarget.value)}
          placeholder="Сообщение..."
          disabled={sending()}
          class="flex-1 px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none disabled:opacity-50"
        />
        <button
          type="submit"
          disabled={sending() || !newMessage().trim()}
          class="px-5 py-2 bg-blue-600 text-white rounded-xl hover:bg-blue-700 disabled:opacity-50 transition font-medium cursor-pointer disabled:cursor-not-allowed"
        >
          {sending() ? "..." : "Отправить"}
        </button>
      </form>
    </div>
  );
};

export default ConversationView;
