import { createSignal } from "solid-js";
import { usePermissions, PERMISSIONS } from "../lib/permissions";
import { api } from "../lib/api";
import { useNavigate } from "@solidjs/router";
import type { Post } from "../lib/types";

// TODO: add photo support
const CreatePost = () => {
  const [name, setName] = createSignal("");
  const [description, setDescription] = createSignal("");
  const [error, setError] = createSignal("");
  const [loading, setLoading] = createSignal(false);
  const { hasPermission } = usePermissions();
  const navigate = useNavigate();

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    setError("");
    setLoading(true);
    const formData = new FormData();
    formData.append("name", name());
    if (description()?.trim()) formData.append("description", description());

    try {
      await api.post<{ posts: Post }>("/posts", formData);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Не удалось создать объявление",
      );
    } finally {
      setLoading(false);
      navigate("/");
    }
  };

  return (
    <>
      {hasPermission(PERMISSIONS.POST_CREATE) && (
        <form
          onSubmit={handleSubmit}
          class="max-w-md mx-auto bg-white rounded-lg shadow-md p-6 space-y-4"
        >
          <input
            type="text"
            value={name()}
            onInput={(e) => setName(e.currentTarget.value)}
            placeholder="Название*"
            required
          />
          <input
            type="text"
            value={description()}
            onInput={(e) => setDescription(e.currentTarget.value)}
            placeholder="Описание"
          />
          {error() && <div class="error">{error()}</div>}
          <button type="submit" disabled={loading()}>
            {loading() ? "Отправка..." : "Отправить"}
          </button>
        </form>
      )}
    </>
  );
};

export default CreatePost;
