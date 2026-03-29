import { createSignal, For, Show, createEffect } from "solid-js";
import { usePermissions, getPermissions } from "../lib/permissions";
import { useAuth } from "../lib/auth";

const Profile = () => {
  const [loading, setLoading] = createSignal(true);
  const [error, setError] = createSignal("");
  const { hasPermission } = usePermissions();
  const { UserReadOwn } = getPermissions();

  const { user } = useAuth();

  createEffect(() => {
    if (user()) setLoading(false);
    else setLoading(true);
  });

  // TODO: move to a separate file, the code is duplicated
  const formatDate = (dateStr: string) => {
    return new Date(dateStr).toLocaleDateString("ru-RU");
  };

  return (
    <>
    {hasPermission(UserReadOwn) &&
    <div class="max-w-4xl mx-auto space-y-6">
      <h1 class="text-2xl font-bold text-center">Профиль</h1>

      <Show when={loading()}>
        <div class="text-center py-8">Загрузка...</div>
      </Show>

      <Show when={error()}>
        <div class="bg-red-100 text-red-700 p-4 rounded-lg">{error()}</div>
      </Show>

      <Show when={!loading() && !error()}>
        <ul>
          <img
            class="w-10 h-10"
            src={`/storage/storage/avatars/${user().hasAvatar ? user().id : "default"}.jpeg`}
            alt="Фото профиля"
          />
          <li>ID пользователя: {user().id}</li>
          <li>Email: {user().email}</li>
          <li>Имя: {user().firstName}</li>
          <li>Фамилия: {user().lastName}</li>
          {user().middleName && <li>Отчество: {user().middleName}</li>}
          <li>Аккаунт создан: {formatDate(user().createdAt)}</li>
          <li>
            Роли:
            <ul>
            <For each={user().roles}>
            {role => <li>{role.name}</li>}
            </For>
            </ul>
          </li>
        </ul>
      </Show>
    </div>
    }
    </>
  );
};

export default Profile;
