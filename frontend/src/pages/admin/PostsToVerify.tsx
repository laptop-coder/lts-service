import {
  createSignal,
  createEffect,
  onMount,
  onCleanup,
  For,
  Show,
} from "solid-js";
import { api } from "../../lib/api";
import { usePermissions, PERMISSIONS } from "../../lib/permissions";
import type { Post } from "../../lib/types";
import PostCardCompact from "../../components/PostCardCompact";

const PostsToVerify = () => {
  const { hasPermission } = usePermissions();

  const [posts, setPosts] = createSignal<Post[]>([]);
  const [loading, setLoading] = createSignal(true);
  const [error, setError] = createSignal("");

  const [page, setPage] = createSignal(0);
  const [hasMore, setHasMore] = createSignal(true);
  let observerRef!: HTMLDivElement;
  let observer: IntersectionObserver;

  // for infinite scroll
  const loadPosts = async () => {
    try {
      if (loading() || !hasMore() || page() === 0) return;
      setLoading(true);
      const data = await api.get<{ posts: Post[] }>(
        `/posts?verified=false&limit=10&offset=${page() * 10}`,
      );
      setPosts([...posts(), ...data.posts]);
      setPage(page() + 1);
      setHasMore(data.posts.length === 10);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Ошибка загрузки объявлений",
      );
    } finally {
      setLoading(false);
    }
  };

  // for first loading and refresh after actions
  const refreshPosts = async () => {
    setPage(0);
    setHasMore(true);
    setLoading(true);
    try {
      const data = await api.get<{ posts: Post[] }>(
        "/posts?verified=false&limit=10&offset=0",
      );
      setPosts(data.posts);
      setPage(1);
      setHasMore(data.posts.length === 10);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Ошибка загрузки объявлений",
      );
    } finally {
      setLoading(false);
    }
  };

  const setupObserver = () => {
    observer?.disconnect();
    observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting && hasMore() && !loading()) {
          loadPosts();
        }
      },
      { threshold: 0.1, rootMargin: "50px" },
    );

    if (observerRef) observer.observe(observerRef);
  };

  createEffect(() => {
    hasMore();
    loading();
    setupObserver();
  });

  onCleanup(() => observer.disconnect());

  onMount(async () => {
    await refreshPosts();
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

        <Show when={loading() && posts().length === 0}>
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
              {(post) => (
                <PostCardCompact post={post} onChange={refreshPosts} />
              )}
            </For>
          </div>
          <div ref={observerRef} class="h-10">
            <Show when={posts().length === 0}>
              <div class="text-center py-16">
                <div class="text-5xl mb-3">📭</div>
                <p class="text-gray-500">Нет объявлений на верификацию</p>
                <p class="text-gray-400 text-sm mt-1">
                  Все объявления проверены
                </p>
              </div>
            </Show>
            <Show when={!hasMore() && posts().length > 0}>
              <div class="text-center text-gray-500 py-8">
                Больше нет объявлений
              </div>
            </Show>
          </div>
        </Show>
      </div>
    )
  );
};

export default PostsToVerify;
