import { createSignal, onMount, For, Show } from "solid-js";
import { A } from "@solidjs/router";
import { ConversationListItem } from "../lib/types";
import { formatDate } from "../lib/utils";
import { conversationApi } from "../lib/api";

const ListOfConversations = () => {
  const [conversations, setConversations] = createSignal<
    ConversationListItem[]
  >([]);
  const [loading, setLoading] = createSignal(true);
  const [error, setError] = createSignal("");

  onMount(async () => {
    try {
      const data = await conversationApi.getListOwn();
      setConversations(data.conversations);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Ошибка загрузки переписок");
    } finally {
      setLoading(false);
    }
  });

  return (
    <div class="max-w-4xl mx-auto p-4 space-y-4">
      <h1 class="text-2xl font-bold text-gray-800">Сообщения</h1>

      <Show when={loading()}>
        <div class="text-center py-8 text-gray-500">Загрузка...</div>
      </Show>

      <Show when={error()}>
        <div class="bg-red-50 text-red-600 p-4 rounded-xl">{error()}</div>
      </Show>

      <Show when={!loading() && conversations().length === 0}>
        <div class="text-center py-16 bg-white rounded-2xl shadow">
          <div class="text-5xl mb-3">💬</div>
          <p class="text-gray-500">Нет сообщений</p>
        </div>
      </Show>

      <Show when={!loading() && conversations().length > 0}>
        <div class="space-y-3">
          <For each={conversations()}>
            {(conv) => (
              <A
                href={`/conversations/${conv.id}`}
                class="block bg-white rounded-xl shadow hover:shadow-md transition p-4"
              >
                <div class="flex items-center gap-4 relative">
                  <img
                    src={`/storage/storage/avatars/${conv.otherUser.hasAvatar ? conv.otherUser.id : "default"}.jpeg`}
                    alt={`Фото профиля пользователя ${conv.otherUser.firstName} ${conv.otherUser.lastName}`}
                    class="w-12 h-12 rounded-full object-cover"
                  />
                  <div class="flex-1 min-w-0">
                    <div class="flex justify-between items-start">
                      <h3 class="font-semibold text-gray-800 truncate">
                        {conv.otherUser.firstName} {conv.otherUser.lastName}
                      </h3>
                      <span class="text-xs text-gray-400">
                        {formatDate(conv.updatedAt)}
                      </span>
                    </div>
                    <p class="text-sm text-gray-500 truncate mt-0.5">
                      {conv.postName}
                    </p>
                    <Show when={conv.lastMessage}>
                      <p class="text-sm text-gray-600 truncate mt-1">
                        {conv.lastMessage}
                      </p>
                    </Show>
                  </div>
                  <Show when={conv.unreadCount > 0}>
                    <div class="bg-blue-600 text-white text-xs font-medium px-2 py-1 rounded-full absolute right-0">
                      {conv.unreadCount}
                    </div>
                  </Show>
                </div>
              </A>
            )}
          </For>
        </div>
      </Show>
    </div>
  );
};

export default ListOfConversations;
