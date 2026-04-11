import { createSignal, createEffect, onMount, For, Show } from "solid-js";
import { api } from "../lib/api";
import type { Post } from "../lib/types";
import PostCardCompact from "../components/PostCardCompact";
import TabsToggle from "../components/TabsToggle";

const PublicPosts = () => {
  const [allPosts, setAllPosts] = createSignal<Post[]>([]);
  const [postsToShow, setPostsToShow] = createSignal<Post[]>([]);
  const [loading, setLoading] = createSignal(true);
  const [error, setError] = createSignal("");

  const loadPosts = async () => {
    try {
      const data = await api.get<{ posts: Post[] }>(
        `/posts/public?thingReturnedToOwner=${activeTab() === tabs[1] ? true : false}`,
      );
      setAllPosts(data.posts);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Ошибка загрузки объявлений",
      );
    } finally {
      setLoading(false);
    }
  };

  createEffect(() => {
    if (activeTab() === tabs[1]) {
      setPostsToShow(
        allPosts().filter((post) => post.thingReturnedToOwner === true),
      );
    } else {
      setPostsToShow(
        allPosts().filter((post) => post.thingReturnedToOwner === false),
      );
    }
  });

  onMount(async () => {
    await loadPosts();
  });

  const tabs = ["Новые", "Закрытые"];
  const [activeTab, setActiveTab] = createSignal(tabs[0]);

  return (
    <div class="max-w-4xl mx-auto space-y-6">
      <h1 class="text-2xl font-bold text-center">Объявления</h1>

      <TabsToggle
        tabs={tabs}
        setActiveTab={setActiveTab}
        tabsHTMLElementId="thing_found_toggle"
        afterChange={loadPosts}
      />

      <Show when={loading()}>
        <div class="text-center py-8">Загрузка...</div>
      </Show>

      <Show when={error()}>
        <div class="bg-red-100 text-red-700 p-4 rounded-lg">{error()}</div>
      </Show>

      <Show when={!loading() && !error()}>
        <div class="space-y-4">
          <For each={postsToShow()}>
            {(post) => <PostCardCompact post={post} onChange={loadPosts} />}
          </For>

          <Show when={postsToShow().length === 0}>
            <div class="text-center text-gray-500 py-8">
              Пока нет объявлений
            </div>
          </Show>
        </div>
      </Show>
    </div>
  );
};

export default PublicPosts;
