import {
  createSignal,
  createEffect,
  createMemo,
  onMount,
  onCleanup,
  Index,
  Show,
} from "solid-js";
import { api } from "../lib/api";
import type { Post } from "../lib/types";
import PostCardCompact from "../components/PostCardCompact";
import TabsToggle from "../components/TabsToggle";
import { usePermissions, ROLES } from "../lib/permissions";
import { useAuth } from "../lib/auth";

const PublicPosts = () => {
  const [posts, setPosts] = createSignal<Post[]>([]);
  const [loading, setLoading] = createSignal(false);
  const [error, setError] = createSignal("");
  const { hasRole } = usePermissions();
  const auth = useAuth();
  const [page, setPage] = createSignal(0);
  const [hasMore, setHasMore] = createSignal(true);
  let observerRef!: HTMLDivElement;
  let observer: IntersectionObserver;

  // for infinite scroll
  const loadPosts = async () => {
    try {
      if (loading() || !hasMore() || page() === 0) return;

      // Save scroll position
      const scrollY = window.scrollY;

      setLoading(true);
      const data = await api.get<{ posts: Post[] }>(
        `/posts/public?author=${ownerTabsActive().query}&thingReturnedToOwner=${statusTabsActive().query}&limit=10&offset=${page() * 10}`,
      );
      setPosts([...posts(), ...data.posts]);
      setPage(page() + 1);
      setHasMore(data.posts.length === 10);

      // Restore scroll position
      requestAnimationFrame(() => {
        window.scrollTo(0, scrollY);
      });
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Ошибка загрузки объявлений",
      );
    } finally {
      setLoading(false);
    }
  };

  // for first loading, for refresh after actions and for refresh when filter changes
  const refreshPosts = async () => {
    setPage(0);
    setHasMore(true);
    setLoading(true);
    try {
      const data = await api.get<{ posts: Post[] }>(
        `/posts/public?author=${ownerTabsActive().query}&thingReturnedToOwner=${statusTabsActive().query}&limit=10&offset=0`,
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

  onMount(async () => {
    await refreshPosts();
  });

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

  // Status tabs
  const statusTabs = [
    { label: "Новые", query: "false" },
    { label: "Закрытые", query: "true" },
  ]; // value in the query answers the question "was the thing in the post returned to owner?"

  // Owner tabs
  const ownerTabs = createMemo(() => {
    const tabs = [
      { label: "Все", query: "all" },
      { label: "Мои", query: "me" },
    ];

    if (!auth.user()?.roles) return tabs;

    if (hasRole(ROLES.TEACHER)) {
      tabs.push({ label: "Мои ученики", query: "students" });
    }
    if (hasRole(ROLES.PARENT)) {
      tabs.push({ label: "Мои дети", query: "children" });
      tabs.push({ label: "Классы детей", query: "children_groups" });
    }
    if (hasRole(ROLES.STUDENT)) {
      tabs.push({ label: "Мои родители", query: "parents" });
      tabs.push({ label: "Мой класс", query: "classmates" });
    }
    return tabs;
  });
  const [ownerTabsActive, setOwnerTabsActive] = createSignal(ownerTabs()[0]);
  const [statusTabsActive, setStatusTabsActive] = createSignal(statusTabs[0]);

  return (
    <div class="max-w-4xl mx-auto space-y-6">
      <h1 class="text-2xl font-bold text-center">Объявления</h1>

      <div class="flex flex-col gap-3">
        <TabsToggle
          tabs={ownerTabs()}
          onChange={(tab) => {
            setOwnerTabsActive(tab);
            refreshPosts();
          }}
          tabsHTMLElementId="owner_tabs_toggle"
        />
        <TabsToggle
          tabs={statusTabs}
          onChange={(tab) => {
            setStatusTabsActive(tab);
            refreshPosts();
          }}
          tabsHTMLElementId="status_tabs_toggle"
        />
      </div>

      <Show when={loading() && posts().length === 0}>
        <div class="text-center py-8">Загрузка...</div>
      </Show>

      <Show when={error()}>
        <div class="bg-red-100 text-red-700 p-4 rounded-lg">{error()}</div>
      </Show>

      <Show when={!loading() && !error()}>
        <div class="space-y-4">
          <Index each={posts()}>
            {(post) => (
              <PostCardCompact post={post()} onChange={refreshPosts} />
            )}
          </Index>
        </div>
      </Show>
      <div ref={observerRef} class="h-10">
        <Show when={posts().length === 0}>
          <div class="text-center text-gray-500 py-8">Пока нет объявлений</div>
        </Show>
        <Show when={!hasMore() && posts().length > 0}>
          <div class="text-center text-gray-500 py-8">
            Больше нет объявлений
          </div>
        </Show>
      </div>
    </div>
  );
};

export default PublicPosts;
