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
import RequestStudentInvite from "./RequestStudentInvite";

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
  const [teacherStudentGroupIds, setTeacherStudentGroupIds] = createStore<
    number[]
  >([]);
  const [studentGroupId, setStudentGroupId] = createSignal<number | null>(null);
  const [staffPositionId, setStaffPositionId] = createSignal<number | null>(
    null,
  );
  const [
    institutionAdministratorPositionId,
    setInstitutionAdministratorPositionId,
  ] = createSignal<number | null>(null);
  const [parentStudentIds, setParentStudentIds] = createStore<string[]>([]);

  const [avatar, setAvatar] = createSignal<File | null>(null);
  const [avatarPreview, setAvatarPreview] = createSignal<string | null>(null);

  const [error, setError] = createSignal("");
  const [loading, setLoading] = createSignal(false);

  const [searchParams] = useSearchParams();
  const auth = useAuth();
  const navigate = useNavigate();

  const inviteToken = searchParams.inviteToken;
  if (typeof inviteToken !== "string") {
    return <RequestStudentInvite />;
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
        setEmailPreloaded(true);
        setEmail(emailData.email);
      }
      // Get roles from the invite token
      const rolesData = await api.get<{ roles: Role[] }>(
        `/tokens/invite/${inviteToken}/roles`,
      );
      setRoleIds(rolesData.roles.map((role) => role.id));
      setRoleNames(rolesData.roles.map((role) => role.name));

      await loadDataForSelect();
    } catch (err) {
      // TODO
    } finally {
      setLoading(false);
    }
  });

  const loadDataForSelect = async () => {
    const [
      roomsData,
      subjectsData,
      groupsData,
      staffPositionData,
      institutionAdministratorPositionData,
    ] = await Promise.all([
      api.get<{ rooms: Room[] }>("/rooms"),
      api.get<{ subjects: Subject[] }>("/subjects"),
      api.get<{ studentGroups: StudentGroup[] }>("/student_groups"),
      api.get<{ staffPositions: StaffPosition[] }>("/staff/positions"),
      api.get<{
        institutionAdministratorPositions: InstitutionAdministratorPosition[];
      }>("/institution_administrators/positions"),
    ]);
    setRooms(roomsData.rooms);
    setSubjects(subjectsData.subjects);
    setStudentGroups(groupsData.studentGroups);
    setStaffPositions(staffPositionData.staffPositions);
    setInstitutionAdministratorPositions(
      institutionAdministratorPositionData.institutionAdministratorPositions,
    );
  };

  const handleAvatarChange = (e: Event) => {
    const input = e.currentTarget as HTMLInputElement;
    const file = input.files?.[0];
    if (file) {
      setAvatar(file);
      const preview = URL.createObjectURL(file);
      setAvatarPreview(preview);
    }
  };

  const removeAvatar = () => {
    setAvatar(null);
    if (avatarPreview()) {
      URL.revokeObjectURL(avatarPreview()!);
      setAvatarPreview(null);
    }
  };

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
    if (avatar()) formData.append("avatar", avatar()!);
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
      teacherStudentGroupIds.forEach((groupId) => {
        formData.append("teacherStudentGroupId", String(groupId));
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

  // Parent students
  const addStudentId = () => {
    setParentStudentIds([...parentStudentIds, ""]);
  };

  const updateStudentId = (index: number, value: string) => {
    setParentStudentIds(index, value);
  };

  const removeStudentId = (index: number) => {
    setParentStudentIds(parentStudentIds.filter((_, i) => i !== index));
  };

  // Student groups where teacher is the advisor
  const addTeacherStudentGroupId = () => {
    setTeacherStudentGroupIds([
      ...teacherStudentGroupIds.map((id) => Number(id)),
      0,
    ]);
  };

  const updateTeacherStudentGroupId = (index: number, value: number) => {
    setTeacherStudentGroupIds(index, value);
  };

  const removeTeacherStudentGroupId = (index: number) => {
    setTeacherStudentGroupIds(
      teacherStudentGroupIds.filter((_, i) => i !== index),
    );
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
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">
                Аватар
              </label>

              <Show when={!avatarPreview()}>
                <label class="flex flex-col items-center justify-center w-full h-32 border-2 border-dashed border-gray-300 rounded-xl cursor-pointer hover:border-blue-500 transition">
                  <div class="flex flex-col items-center justify-center pt-5 pb-6">
                    <svg
                      class="w-8 h-8 text-gray-400 mb-2"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
                      ></path>
                    </svg>
                    <p class="text-sm text-gray-500">
                      Нажмите для загрузки аватара
                    </p>
                    <p class="text-xs text-gray-400 mt-1">
                      JPEG, PNG, WebP, GIF (макс. 10MB)
                    </p>
                  </div>
                  <input
                    type="file"
                    accept="image/jpeg,image/png,image/webp,image/gif"
                    onChange={handleAvatarChange}
                    class="hidden"
                  />
                </label>
              </Show>

              <Show when={avatarPreview()}>
                <div class="relative inline-block">
                  <img
                    src={avatarPreview()!}
                    alt="Preview"
                    class="w-24 h-24 rounded-full object-cover border-4 border-blue-100"
                  />
                  <button
                    type="button"
                    onClick={removeAvatar}
                    class="absolute -top-2 -right-2 p-1 bg-red-500 text-white rounded-full hover:bg-red-600 transition"
                  >
                    <svg
                      class="w-4 h-4"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M6 18L18 6M6 6l12 12"
                      ></path>
                    </svg>
                  </button>
                </div>
              </Show>
            </div>

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
                <div class="space-y-3 border-t border-gray-200 pt-4">
                  <h3 class="font-medium text-gray-800">
                    Классное руководство/менторство
                  </h3>

                  <Index each={teacherStudentGroupIds}>
                    {(groupId, index) => (
                      <div class="flex gap-2">
                        <select
                          value={groupId() || ""}
                          onChange={(e) =>
                            updateTeacherStudentGroupId(
                              index,
                              Number(e.target.value),
                            )
                          }
                          class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition bg-white cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
                        >
                          <option value="">Выберите группу</option>
                          <For
                            each={studentGroups().filter(
                              (group) =>
                                !teacherStudentGroupIds
                                  .filter((_, i) => i !== index)
                                  .includes(group.id),
                            )}
                          >
                            {(group) => (
                              <option value={group.id}>{group.name}</option>
                            )}
                          </For>
                        </select>
                        <button
                          type="button"
                          onClick={() => removeTeacherStudentGroupId(index)}
                          class="px-4 py-2 bg-red-700 text-white rounded-xl hover:bg-red-800 transition cursor-pointer"
                        >
                          Удалить
                        </button>
                      </div>
                    )}
                  </Index>

                  <button
                    type="button"
                    onClick={addTeacherStudentGroupId}
                    class="w-full px-4 py-2 bg-gray-100 text-gray-700 rounded-xl hover:bg-gray-200 transition font-medium cursor-pointer"
                  >
                    + Добавить группу
                  </button>
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
                  class="w-full px-4 py-2 bg-gray-100 text-gray-700 rounded-xl hover:bg-gray-200 transition font-medium cursor-pointer"
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
