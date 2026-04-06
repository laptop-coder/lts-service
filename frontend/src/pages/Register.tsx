import { createSignal, onMount, Show, For, Index } from "solid-js";
import { createStore } from "solid-js/store";
import { useAuth } from "../lib/auth";
import { api } from "../lib/api";
import { ROLES_TO_DISPLAY } from "../lib/permissions";
import { useNavigate, useSearchParams } from "@solidjs/router";
import type {
  Role,
  Room,
  Subject,
  StudentGroup,
  StaffPosition,
  InstitutionAdministratorPosition,
} from "../lib/types";

// TODO: add avatar support
const Register = () => {
  // Data about new user
  const [email, setEmail] = createSignal("");
  const [emailPreloaded, setEmailPreloaded] = createSignal(false);
  const [password, setPassword] = createSignal("");
  const [firstName, setFirstName] = createSignal("");
  const [middleName, setMiddleName] = createSignal("");
  const [lastName, setLastName] = createSignal("");
  // Special fields
  const [teacherClassroomId, setTeacherClassroomId] = createSignal<
    number | null
  >(null);
  const [teacherSubjectIds, setTeacherSubjectIds] = createSignal<number[]>([]);
  const [studentGroupId, setStudentGroupId] = createSignal<number | null>(null);
  const [staffPositionId, setStaffPositionId] = createSignal<number | null>(
    null,
  );
  const [
    institutionAdministratorPositionId,
    setInstitutionAdministratorPositionId,
  ] = createSignal<number | null>(null);
  const [parentStudentIds, setParentStudentIds] = createStore<string[]>([]);

  const [error, setError] = createSignal("");
  const [loading, setLoading] = createSignal(false);

  const [searchParams] = useSearchParams();
  const auth = useAuth();
  const navigate = useNavigate();

  const inviteToken = searchParams.inviteToken;
  if (typeof inviteToken !== "string") {
    return <div>Для регистрации нужен инвайт-токен</div>;
  }
  // Roles
  const [roleIds, setRoleIds] = createSignal<number[]>([]);
  const [roleNames, setRoleNames] = createSignal<string[]>([]);
  // Rooms
  const [rooms, setRooms] = createSignal<Room[]>([]);
  // Subjects
  const [subjects, setSubjects] = createSignal<Subject[]>([]);
  // Student groups
  const [studentGroups, setStudentGroups] = createSignal<StudentGroup[]>([]);
  // Staff positions
  const [staffPositions, setStaffPositions] = createSignal<StaffPosition[]>([]);
  // Institution administrators positions
  const [
    institutionAdministratorPositions,
    setInstitutionAdministratorPositions,
  ] = createSignal<InstitutionAdministratorPosition[]>([]);

  onMount(async () => {
    try {
      // Try to get email from the invite token
      const emailData = await api.get<{ email: string }>(
        `/tokens/invite/${inviteToken}/email`,
      );
      if (emailData.email) {
        setEmailPreloaded(true)
        setEmail(emailData.email)
      }
      // Get roles from the invite token
      const rolesData = await api.get<{ roles: Role[] }>(
        `/tokens/invite/${inviteToken}/roles`,
      );
      setRoleIds(rolesData.roles.map((role) => role.id));
      setRoleNames(rolesData.roles.map((role) => role.name));
      // Get all rooms
      const roomsData = await api.get<{ rooms: Room[] }>("/rooms");
      setRooms(roomsData.rooms);
      // Get all subjects
      const subjectsData = await api.get<{ subjects: Subject[] }>("/subjects");
      setSubjects(subjectsData.subjects);
      // Get all student groups
      const studentGroupsData = await api.get<{
        studentGroups: StudentGroup[];
      }>("/student_groups");
      setStudentGroups(studentGroupsData.studentGroups);
      // Get all staff positions
      const staffPositionsData = await api.get<{
        staffPositions: StaffPosition[];
      }>("/staff/positions");
      setStaffPositions(staffPositionsData.staffPositions);
      // Get all institution administrator positions
      const institutionAdministratorPositionsData = await api.get<{
        institutionAdministratorPositions: InstitutionAdministratorPosition[];
      }>("/institution_administrators/positions");
      setInstitutionAdministratorPositions(
        institutionAdministratorPositionsData.institutionAdministratorPositions,
      );
    } catch (err) {
      // TODO
    } finally {
      setLoading(false);
    }
  });

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    const formData = new FormData();
    if (!emailPreloaded()) formData.append("email", email());
    formData.append("password", password());
    formData.append("firstName", firstName());
    if (middleName()?.trim()) formData.append("middleName", middleName());
    formData.append("lastName", lastName());
    formData.append("inviteToken", inviteToken);
    if (roleIds().includes(3)) {
      if (!institutionAdministratorPositionId()) {
        setError("Выберите должность администрации ОУ");
        return;
      }
      formData.append(
        "institutionAdministratorPositionId",
        String(institutionAdministratorPositionId()),
      );
    }
    if (roleIds().includes(4)) {
      if (!staffPositionId()) {
        setError("Выберите должность сотрудника ОУ");
        return;
      }
      formData.append("staffPositionId", String(staffPositionId()));
    }
    if (roleIds().includes(5)) {
      if (!teacherClassroomId()) {
        setError("Выберите аудиторию/классный кабинет");
        return;
      }
      formData.append("teacherClassroomId", String(teacherClassroomId()));
      if (teacherSubjectIds().length === 0) {
        setError("Выберите предметы");
        return;
      }
      teacherSubjectIds().forEach((subjectId) => {
        formData.append("teacherSubjectId", String(subjectId));
      });
    }
    if (roleIds().includes(6)) {
      parentStudentIds.forEach((studentId) => {
        if (studentId.trim()) {
          formData.append("parentStudentId", studentId);
        }
      });
    }
    if (roleIds().includes(7)) {
      if (!studentGroupId()) {
        setError("Выберите класс/учебную группу");
        return;
      }
      formData.append("studentGroupId", String(studentGroupId()));
    }

    try {
      await auth.register(formData);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Ошибка регистрации");
    } finally {
      setLoading(false);
      navigate("/");
    }
  };

  const addStudentId = () => {
    setParentStudentIds([...parentStudentIds, ""]);
  };

  const updateStudentId = (index: number, value: string) => {
    setParentStudentIds(index, value);
  };

  const removeStudentId = (index: number) => {
    setParentStudentIds(parentStudentIds.filter((_, i) => i !== index));
  };

  return (
    <div class="min-h-screen py-8 px-4">
      <div class="max-w-2xl mx-auto">
        <div class="text-center mb-8">
          <h1 class="text-3xl font-bold text-gray-800">Регистрация</h1>
          <p class="text-gray-500 mt-2">
            Роли:{" "}
            <span class="font-medium text-blue-600">
              {ROLES_TO_DISPLAY.filter((role) =>
                roleNames().includes(role.name),
              )
                .map((role) => role.displayName)
                .join(", ")}
            </span>
          </p>
        </div>

        {loading() ? (
          <div class="text-center py-12 text-gray-500">Загрузка...</div>
        ) : (
          <form
            onSubmit={handleSubmit}
            class="bg-white rounded-2xl shadow-lg p-6 space-y-5"
          >
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">
                  Фамилия *
                </label>
                <input
                  type="text"
                  value={lastName()}
                  placeholder="Иванов"
                  onInput={(e) => setLastName(e.currentTarget.value)}
                  class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition"
                  required
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">
                  Имя *
                </label>
                <input
                  type="text"
                  value={firstName()}
                  placeholder="Иван"
                  onInput={(e) => setFirstName(e.currentTarget.value)}
                  class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition"
                  required
                />
              </div>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">
                Отчество
              </label>
              <input
                type="text"
                value={middleName()}
                placeholder="Иванович"
                onInput={(e) => setMiddleName(e.currentTarget.value)}
                class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition"
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">
                Email *
              </label>
              <input
              disabled={emailPreloaded()}
                type="email"
                value={email()}
                placeholder="email@example.ru"
                onInput={(e) => setEmail(e.currentTarget.value)}
                class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition disabled:cursor-not-allowed"
                required
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">
                Пароль *
              </label>
              <input
                type="password"
                value={password()}
                placeholder="••••••••"
                onInput={(e) => setPassword(e.currentTarget.value)}
                class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition"
                required
              />
            </div>

            <Show when={roleIds().includes(3)}>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">
                  Должность администрации ОУ *
                </label>
                <select
                  value={institutionAdministratorPositionId() || ""}
                  onChange={(e) =>
                    setInstitutionAdministratorPositionId(
                      Number(e.currentTarget.value),
                    )
                  }
                  class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition bg-white cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  <option value="">Выберите должность</option>
                  <For each={institutionAdministratorPositions()}>
                    {(position) => (
                      <option value={position.id}>{position.name}</option>
                    )}
                  </For>
                </select>
              </div>
            </Show>

            <Show when={roleIds().includes(4)}>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">
                  Должность сотрудника ОУ *
                </label>
                <select
                  value={staffPositionId() || ""}
                  onChange={(e) =>
                    setStaffPositionId(Number(e.currentTarget.value))
                  }
                  class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition bg-white cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  <option value="">Выберите должность</option>
                  <For each={staffPositions()}>
                    {(position) => (
                      <option value={position.id}>{position.name}</option>
                    )}
                  </For>
                </select>
              </div>
            </Show>

            <Show when={roleIds().includes(5)}>
              <div class="space-y-4 border-t border-gray-200 pt-4">
                <h3 class="font-medium text-gray-800">Данные учителя</h3>

                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">
                    Классный кабинет *
                  </label>
                  <select
                    value={teacherClassroomId() || ""}
                    onChange={(e) =>
                      setTeacherClassroomId(Number(e.currentTarget.value))
                    }
                    class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition bg-white cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    <option value="">Выберите кабинет</option>
                    <For each={rooms()}>
                      {(room) => <option value={room.id}>{room.name}</option>}
                    </For>
                  </select>
                </div>

                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">
                    Предметы *
                  </label>
                  <div class="grid grid-cols-2 gap-2 max-h-48 overflow-y-auto p-2 border border-gray-200 rounded-xl">
                    <For each={subjects()}>
                      {(subject) => (
                        <label class="flex items-center gap-2 p-2 hover:bg-gray-50 rounded-lg cursor-pointer transition">
                          <input
                            type="checkbox"
                            checked={teacherSubjectIds().includes(subject.id)}
                            onChange={() => {
                              if (teacherSubjectIds().includes(subject.id)) {
                                setTeacherSubjectIds((prev) =>
                                  prev.filter((id) => id !== subject.id),
                                );
                              } else {
                                setTeacherSubjectIds((prev) => [
                                  ...prev,
                                  subject.id,
                                ]);
                              }
                            }}
                            class="w-4 h-4 text-blue-600 rounded focus:ring-blue-500"
                          />
                          <span class="text-gray-700 text-sm">
                            {subject.name}
                          </span>
                        </label>
                      )}
                    </For>
                  </div>
                </div>
              </div>
            </Show>

            <Show when={roleIds().includes(7)}>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">
                  Класс/учебная группа *
                </label>
                <select
                  value={studentGroupId() || ""}
                  onChange={(e) =>
                    setStudentGroupId(Number(e.currentTarget.value))
                  }
                  class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition bg-white cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  <option value="">Выберите группу</option>
                  <For each={studentGroups()}>
                    {(group) => <option value={group.id}>{group.name}</option>}
                  </For>
                </select>
              </div>
            </Show>

            <Show when={roleIds().includes(6)}>
              <div class="space-y-3 border-t border-gray-200 pt-4">
                <h3 class="font-medium text-gray-800">Привязка учеников</h3>

                <Index each={parentStudentIds}>
                  {(studentId, index) => (
                    <div class="flex gap-2">
                      <input
                        type="text"
                        value={studentId()}
                        onInput={(e) => updateStudentId(index, e.target.value)}
                        placeholder={`ID ученика ${index + 1}`}
                        class="flex-1 px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition"
                      />
                      <button
                        type="button"
                        onClick={() => removeStudentId(index)}
                        class="px-4 py-2 bg-red-700 text-white rounded-xl hover:bg-red-800 transition cursor-pointer"
                      >
                        Удалить
                      </button>
                    </div>
                  )}
                </Index>

                <button
                  type="button"
                  onClick={addStudentId}
                  class="w-full px-4 py-2 bg-gray-100 text-gray-700 rounded-xl hover:bg-gray-200 transition font-medium"
                >
                  + Добавить ученика
                </button>
              </div>
            </Show>

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
                {loading() ? "Регистрация..." : "Зарегистрироваться"}
              </button>
            </div>
          </form>
        )}
      </div>
    </div>
  );
};

export default Register;
