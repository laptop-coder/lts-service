import { createSignal, onMount, For, Show } from "solid-js";
import { api } from "../../lib/api";
import { usePermissions, PERMISSIONS } from "../../lib/permissions";
import type { Post } from "../../lib/types";
import PostCardCompact from "../../components/PostCardCompact";

const PostsToVerify = () => {
  const { hasPermission } = usePermissions();

  const [posts, setPosts] = createSignal<Post[]>([]);
  const [loading, setLoading] = createSignal(true);
  const [error, setError] = createSignal("");

  const loadPosts = async () => {
    try {
      const data = await api.get<{ posts: Post[] }>("/posts?verified=false");
      setPosts(data.posts);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Ошибка загрузки объявлений",
      );
    } finally {
      setLoading(false);
    }
  };

  onMount(() => {
    loadPosts();
  });

  return (
    hasPermission(PERMISSIONS.POST_READ_ANY) && (
      <div class="space-y-6 p-4">
        <div class="mb-6">
          <h1 class="text-3xl font-bold text-gray-800">
            Верификация объявлений
          </h1>
          <p class="text-gray-500 mt-1">
            Проверьте и подтвердите объявления пользователей
          </p>
        </div>

        <Show when={loading()}>
          <div class="flex justify-center items-center py-16">
            <div class="text-gray-500">Загрузка...</div>
          </div>
        </Show>

        <Show when={error()}>
          <div class="bg-red-50 border border-red-200 text-red-600 p-4 rounded-xl">
            {error()}
          </div>
        </Show>

        <Show when={!loading() && !error()}>
          <div class="space-y-4">
            <For each={posts()}>
              {(post) => <PostCardCompact post={post} onChange={loadPosts} />}
            </For>

            <Show when={posts().length === 0}>
              <div class="text-center py-16">
                <div class="text-5xl mb-3">📭</div>
                <p class="text-gray-500">Нет объявлений на верификацию</p>
                <p class="text-gray-400 text-sm mt-1">
                  Все объявления проверены
                </p>
              </div>
            </Show>
          </div>
        </Show>
      </div>
    )
  );
};

export default PostsToVerify;
