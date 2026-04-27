import { A } from "@solidjs/router";
import { Show } from "solid-js";
import {
  Search,
  MessageCircleQuestionMark,
  MessageSquareText,
  Shield,
  GraduationCap,
  BookOpen,
  Users,
  Building2,
  UserCog,
} from "lucide-solid";
import { useAuth } from "../lib/auth";

const About = () => {
  const auth = useAuth();
  return (
    <div class="max-w-5xl mx-auto px-4 py-8 md:py-16 space-y-12">
      <div class="text-center space-y-4">
        <h1 class="text-3xl md:text-5xl font-bold text-gray-800">
          LostThingsSearch
        </h1>
        <p class="text-lg md:text-xl text-gray-500 max-w-2xl mx-auto">
          Сервис поиска потерянных вещей для образовательных учреждений
        </p>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div class="bg-white rounded-2xl p-6 shadow-sm">
          <h2 class="text-lg font-semibold text-gray-800 mb-2 flex gap-1">
            <Search /> Нашли что-то чужое?
          </h2>
          <p class="text-gray-600 text-sm">
            Сфотографируйте находку и создайте объявление. Укажите, где и когда вы её нашли, чтобы владелец мог быстро найти свою вещь.
          </p>
        </div>
        <div class="bg-white rounded-2xl p-6 shadow-sm">
          <h2 class="text-lg font-semibold text-gray-800 mb-2 flex gap-1">
            <MessageCircleQuestionMark /> Потеряли своё?
          </h2>
          <p class="text-gray-600 text-sm">
            Проверьте ленту объявлений — возможно, кто-то уже ищет хозяина вашей вещи. Используйте фильтры, чтобы ускорить поиск.
          </p>
        </div>
        <div class="bg-white rounded-2xl p-6 shadow-sm">
          <h2 class="text-lg font-semibold text-gray-800 mb-2 flex gap-1">
            <MessageSquareText /> Свяжитесь с автором
          </h2>
          <p class="text-gray-600 text-sm">
            Напишите автору объявления через встроенный чат, чтобы договориться о возврате, не раскрывая личных контактов.
          </p>
        </div>
        <div class="bg-white rounded-2xl p-6 shadow-sm">
          <h2 class="text-lg font-semibold text-gray-800 mb-2 flex gap-1">
            <Shield /> Безопасность
          </h2>
          <p class="text-gray-600 text-sm">
            Все объявления проверяются администраторами сервиса.
          </p>
        </div>
      </div>

      <div class="space-y-6">
        <h2 class="text-2xl font-bold text-gray-800 text-center">
          Для кого этот сервис
        </h2>
        <div class="grid grid-cols-1 md:grid-cols-5 gap-4">
          <div class="bg-white rounded-2xl p-4 shadow-sm text-center">
            <GraduationCap class="text-3xl mb-2 inline-flex" />
            <h3 class="font-semibold text-gray-800">Ученики</h3>
          </div>
          <div class="bg-white rounded-2xl p-4 shadow-sm text-center">
            <BookOpen class="text-3xl mb-2 inline-flex" />
            <h3 class="font-semibold text-gray-800">Учителя</h3>
          </div>
          <div class="bg-white rounded-2xl p-4 shadow-sm text-center">
            <Users class="text-3xl mb-2 inline-flex" />
            <h3 class="font-semibold text-gray-800">Родители</h3>
          </div>
          <div class="bg-white rounded-2xl p-4 shadow-sm text-center">
            <UserCog class="text-3xl mb-2 inline-flex" />
            <h3 class="font-semibold text-gray-800">Сотрудники</h3>
          </div>
          <div class="bg-white rounded-2xl p-4 shadow-sm text-center">
            <Building2 class="text-3xl mb-2 inline-flex" />
            <h3 class="font-semibold text-gray-800">Администрация</h3>
          </div>
        </div>
      </div>

      <Show when={!auth.user()}>
        <div class="text-center space-y-4 bg-white rounded-2xl p-8 shadow-sm">
          <h2 class="text-2xl font-bold text-gray-800">Присоединяйтесь!</h2>
          <A
            href="/register"
            class="w-36 h-10 bg-blue-600 text-white rounded-xl hover:bg-blue-700 transition font-medium inline-flex items-center justify-center"
          >
            Регистрация
          </A>
        </div>
      </Show>
    </div>
  );
};

export default About;
