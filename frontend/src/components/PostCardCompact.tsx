import { Show, createSignal } from "solid-js";
import type { Post } from "../lib/types";
import { usePermissions, getPermissions } from "../lib/permissions";
import { api } from "../lib/api";

interface Props {
  post: Post;
}

const PostCardCompact = (props: Props) => {
  const { post } = props;
  const { hasPermission } = usePermissions();
  const { PostVerify, PostDeleteAny } = getPermissions();
  const [loading, setLoading] = createSignal(false);
  const [error, setError] = createSignal("");

  const verifyPost = async () => {
    try {
      await api.patch<{ posts: Post[] }>(`/posts/${post.id}/verify`);
      setLoading(true)
    } catch (err) {
      setError(err instanceof Error ? err.message : "Не удалось верифицировать объявление");
    } finally {
      setLoading(false);
    }
  };

  const deletePost = async () => {
    try {
      await api.delete<{}>(`/posts/${post.id}`);
      setLoading(true)
    } catch (err) {
      setError(err instanceof Error ? err.message : "Не удалось удалить объявление");
    } finally {
      setLoading(false);
    }
  };


  const formatDate = (dateStr: string) => {
    return new Date(dateStr).toLocaleDateString("ru-RU");
  };

  return (
    <div class="bg-white rounded-lg shadow-md hover:shadow-lg transition overflow-hidden">
      <div class="flex flex-col sm:flex-row">
        <img class='w-10 h-10' src={`/storage/storage/avatars/${post.author.hasAvatar ?  post.author.id : "default"}.jpeg`} alt="Фото профиля" /> {post.name}
        Последнее изменение: {formatDate(post.updatedAt)}
        {post.author.firstName}
        {post.author.lastName}
        <Show when={post.description}>{post.description}</Show>
        {hasPermission(PostVerify) && !post.verified && <button onClick={verifyPost}>Верифицировать</button>}
        {hasPermission(PostDeleteAny) && <button onClick={deletePost}>Удалить</button>}
      </div>
    </div>
  );
};

export default PostCardCompact;
