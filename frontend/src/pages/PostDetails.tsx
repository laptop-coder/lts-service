import { createSignal, Show, onMount } from "solid-js";
import { usePermissions, PERMISSIONS } from "../lib/permissions";
import { api } from "../lib/api";
import { useAuth } from "../lib/auth";
import { useParams } from "@solidjs/router";
import type { Post } from "../lib/types";
import PostCardDetailed from "../components/PostCardDetailed";

const PostDetails = () => {
  const params = useParams();
  const [error, setError] = createSignal("");
  const [loading, setLoading] = createSignal(true);
  const { hasPermission } = usePermissions();
  const auth = useAuth();
  const [post, setPost] = createSignal<Post | null>(null);

  const loadPost = async () => {
    try {
      const data = await api.get<{ post: Post }>(`/posts/${params.id}`);
      setPost(data.post);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Не удалось загрузить объявление",
      );
    } finally {
      setLoading(false);
    }
  };

  onMount(() => {
    loadPost();
  });

  return (
    <div class="max-w-4xl mx-auto px-4 py-6">
            <h1 class="text-2xl font-bold text-gray-800 text-center mb-6">
            Информация об объявлении
            </h1>
      <Show when={loading()}>
        <div class="text-center py-8">Загрузка...</div>;
      </Show>
      <Show when={error()}>
        <div class="bg-red-50 text-red-600 p-3 rounded-xl text-sm border border-red-200">
          {error()}
        </div>
      </Show>
      <Show when={!loading() && post()}>
        <Show
          when={
            post()!.verified ||
            hasPermission(PERMISSIONS.POST_READ_ANY) ||
            (hasPermission(PERMISSIONS.POST_READ_OWN) &&
              post()!.author.id === auth.user()?.id)
          }
        >
          <PostCardDetailed post={post()!} onChange={loadPost} />
        </Show>
      </Show>
    </div>
  );
};

export default PostDetails;
