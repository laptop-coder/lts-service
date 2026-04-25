import { createSignal, createEffect } from "solid-js";
import { useNavigate } from "@solidjs/router";
import { api } from "../lib/api";
import { Mail, ChevronLeft } from "lucide-solid";
import { ROLES } from "../lib/permissions";

const RequestStudentOrParentInvite = () => {
  const [email, setEmail] = createSignal("");
  const [error, setError] = createSignal("");
  const [loading, setLoading] = createSignal(false);
  const [role, setRole] = createSignal<string | null>(null);
  const [done, setDone] = createSignal(false); // email was sent
  const navigate = useNavigate();

  let emailInputRef: HTMLInputElement | undefined;
  const focusEmailInput = () => {
    if (emailInputRef) {
      emailInputRef.focus();
    }
  };

  createEffect(() => {
    role() !== null;
    focusEmailInput();
  });

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    if (!email().trim()) return;
    setError("");
    setLoading(true);

    const formData = new URLSearchParams();
    formData.append("email", email());

    try {
      switch (role()) {
        case ROLES.STUDENT:
          await api.post("/invite/request/student", formData);
          break;
        case ROLES.PARENT:
          await api.post("/invite/request/parent", formData);
          break;
        default:
          setError("Роль не выбрана");
      }
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : "Ошибка запроса пригласительной ссылки",
      );
    } finally {
      setLoading(false);
      setDone(true);
    }
  };

  return (
    <div class="min-h-[80vh] flex items-center justify-center px-4">
      <div class="w-full max-w-md">
        <div class="text-center mb-8">
          <h1 class="text-3xl font-bold text-gray-800">
            Создание аккаунта
            {role() === ROLES.STUDENT && " ученика"}
            {role() === ROLES.PARENT && " родителя"}
          </h1>
          <p class="text-gray-500 mt-2">Запрос пригласительной ссылки</p>
        </div>

        {loading() ? (
          <div class="text-center py-12 text-gray-500">Загрузка...</div>
        ) : !role() ? (
          <div class="bg-white rounded-2xl shadow-lg p-6 flex flex-col gap-4">
            <button
              onClick={() => navigate("/login")}
              type="button"
              class="text-gray-500 hover:text-gray-700 cursor-pointer flex flex-row"
            >
              <ChevronLeft /> Назад
            </button>
            <div class="flex flex-row gap-4">
              <button
                type="button"
                disabled={loading()}
                class="px-4 py-2 bg-blue-600 text-white rounded-xl hover:bg-blue-700 transition font-medium disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer w-full h-full"
                onClick={() => setRole(ROLES.STUDENT)}
              >
                Я ученик
              </button>
              <button
                type="button"
                disabled={loading()}
                class="px-4 py-2 bg-blue-600 text-white rounded-xl hover:bg-blue-700 transition font-medium disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer w-full h-full"
                onClick={() => setRole(ROLES.PARENT)}
              >
                Я родитель
              </button>
            </div>
          </div>
        ) : done() ? (
          <div class="bg-green-100 border border-green-400 text-green-800 px-4 py-3 rounded-xl flex items-center gap-3">
            <Mail />
            <span>
              Письмо отправлено на{" "}
              <span class="font-semibold underline">{email()}</span>! Перейдите
              по ссылке в письме для завершения регистрации.
            </span>
          </div>
        ) : (
          <form
            onSubmit={handleSubmit}
            class="bg-white rounded-2xl shadow-lg p-6 space-y-5"
          >
            <button
              type="button"
              onClick={() => setRole(null)}
              class="text-gray-500 hover:text-gray-700 cursor-pointer flex flex-row"
            >
              <ChevronLeft /> Назад
            </button>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">
                Email *
              </label>
              <input
                ref={emailInputRef}
                type="email"
                value={email()}
                placeholder={`${role() === ROLES.STUDENT ? "student" : "parent"}@example.ru`}
                onInput={(e) => setEmail(e.currentTarget.value)}
                class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition disabled:cursor-not-allowed"
                required
              />
            </div>

            {error() && (
              <div class="bg-red-50 text-red-600 p-3 rounded-xl text-sm border border-red-200">
                {error()}
              </div>
            )}

            <div class="flex gap-3 pt-2">
              <button
                type="submit"
                disabled={loading()}
                class="flex-1 px-4 py-2 bg-blue-600 text-white rounded-xl hover:bg-blue-700 transition font-medium disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
              >
                {loading() ? "Отправка..." : "Отправить"}
              </button>
            </div>
          </form>
        )}
      </div>
    </div>
  );
};

export default RequestStudentOrParentInvite;
