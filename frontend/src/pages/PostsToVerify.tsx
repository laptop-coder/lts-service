import { createSignal, onMount, For, Show } from "solid-js";
import { api } from "../lib/api";
import { usePermissions, getPermissions } from "../lib/permissions";
import type { Post } from "../lib/types";
import PostCardCompact from "../components/PostCardCompact";

const PostsToVerify = () => {
  const { hasPermission } = usePermissions();
  const { PostReadAny } = getPermissions();

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
    <>
      {hasPermission(PostReadAny) && (
        <div class="max-w-4xl mx-auto space-y-6">
          <h1 class="text-2xl font-bold text-center">Верификация объявлений</h1>

          <Show when={loading()}>
            <div class="text-center py-8">Загрузка...</div>
          </Show>

          <Show when={error()}>
            <div class="bg-red-100 text-red-700 p-4 rounded-lg">{error()}</div>
          </Show>

          <Show when={!loading() && !error()}>
            <div class="space-y-4">
              <For each={posts()}>
                {(post) => <PostCardCompact post={post} />}
              </For>

              <Show when={posts().length === 0}>
                <div class="text-center text-gray-500 py-8">
                  Пока нет объявлений
                </div>
              </Show>
            </div>
          </Show>
        </div>
      )}
    </>
  );
};

export default PostsToVerify;
