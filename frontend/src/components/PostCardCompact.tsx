import { Show } from "solid-js";
import type { Post } from "../lib/types";

interface Props {
  post: Post;
}

const PostCardCompact = (props: Props) => {
  const { post } = props;

  const formatDate = (dateStr: string) => {
    return new Date(dateStr).toLocaleDateString("ru-RU");
  };

  return (
    <div class="bg-white rounded-lg shadow-md hover:shadow-lg transition overflow-hidden">
      <div class="flex flex-col sm:flex-row">
        <img class='w-10 h-10' src={`/storage/storage/avatars/${post.author.hasAvatar ? "default" : post.author.id}.jpeg`} alt="Фото профиля" /> {post.name}
        Последнее изменение: {formatDate(post.updatedAt)}
        {post.author.firstName}
        {post.author.lastName}
        <Show when={post.description}>{post.description}</Show>
      </div>
    </div>
  );
};

export default PostCardCompact;
