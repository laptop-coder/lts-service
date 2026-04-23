import {
  createSignal,
  createEffect,
  createMemo,
  onMount,
  onCleanup,
  For,
  Show,
  on,
} from "solid-js";
import { api } from "../lib/api";
import type {
  Post,
  StudentGroup,
  Student,
  StudentGroupStudent,
  Parent,
} from "../lib/types";
import PostCardCompact from "../components/PostCardCompact";
import TabsToggle from "../components/TabsToggle";
import { usePermissions, ROLES } from "../lib/permissions";
import { useAuth } from "../lib/auth";

const PublicPosts = () => {
  const [allPosts, setAllPosts] = createSignal<Post[]>([]);
  const [loading, setLoading] = createSignal(false);
  const [error, setError] = createSignal("");
  const { hasRole } = usePermissions();
  const auth = useAuth();
  const [page, setPage] = createSignal(0);
  const [hasMore, setHasMore] = createSignal(true);
  let observerRef!: HTMLDivElement;
  let observer: IntersectionObserver;

  const loadPosts = async () => {
    try {
      if (loading() || !hasMore()) return;
      setLoading(true);
      const data = await api.get<{ posts: Post[] }>(
        `/posts/public?limit=10&offset=${page() * 10}`,
      );
      setAllPosts([...allPosts(), ...data.posts]);
      setPage(page() + 1);
      setHasMore(data.posts.length === 10);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Ошибка загрузки объявлений",
      );
    } finally {
      setLoading(false);
    }
  };

  const [studentGroupsWhereAdvisor, setStudentGroupsWhereAdvisor] =
    createSignal<StudentGroup[]>([]);
  const [children, setChildren] = createSignal<Student[]>([]);
  const [childrenStudentGroups, setChildrenStudentGroups] = createSignal<
    StudentGroup[]
  >([]);
  const [studentParents, setStudentParents] = createSignal<Parent[]>([]);
  const [classmates, setClassmates] = createSignal<StudentGroupStudent[]>([]);

  const loadRoleData = async () => {
    try {
      if (hasRole(ROLES.TEACHER)) await loadTeacherData();
      if (hasRole(ROLES.PARENT)) await loadParentData();
      if (hasRole(ROLES.STUDENT)) await loadStudentData();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Ошибка загрузки данных");
    }
  };

  const loadTeacherData = async () => {
    const data = await api.get<{ studentGroups: StudentGroup[] }>(
      "/teachers/me/student_groups",
    );
    setStudentGroupsWhereAdvisor(data.studentGroups);
  };

  const loadParentData = async () => {
    const children = await api.get<{ students: Student[] }>(
      "/parents/me/students",
    );
    setChildren(children.students);

    const groupIds = [
      ...new Set(children.students.map((student) => student.studentGroup.id)),
    ];
    const groups = await Promise.all(
      groupIds.map((id) =>
        api.get<{ studentGroup: StudentGroup }>(`/student_groups/${id}`),
      ),
    );
    setChildrenStudentGroups(groups.map((r) => r.studentGroup));
  };

  const loadStudentData = async () => {
    const [parents, group] = await Promise.all([
      api.get<{ parents: Parent[] }>("/students/me/parents"),
      api.get<{ studentGroup: StudentGroup }>("/students/me/student_group"),
    ]);

    setStudentParents(parents.parents);
    setClassmates(group.studentGroup.students);
  };

  createEffect(
    on(
      () => auth.user(),
      (user) => user && loadRoleData(),
    ),
  );

  const resetAndLoad = async () => {
    setPage(0);
    setHasMore(true);
    try {
      const data = await api.get<{ posts: Post[] }>(
        "/posts/public?limit=10&offset=0",
      );
      setAllPosts(data.posts);
      setPage(1);
      setHasMore(data.posts.length === 10);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Ошибка загрузки объявлений",
      );
    }
  };

  onMount(async () => {
    await loadPosts();
  });

  const setupObserver = () => {
    observer?.disconnect();
    observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting && hasMore()) {
          loadPosts();
        }
      },
      { threshold: 0.1 },
    );

    if (observerRef) observer.observe(observerRef);
  };

  createEffect(() => {
    hasMore();
    loading();
    setupObserver();
  });

  onCleanup(() => observer.disconnect());

  // Status tabs
  const statusTabs = ["Новые", "Закрытые"];

  // Owner tabs
  const ownerTabs = createMemo(() => {
    const tabs = ["Все", "Мои"];

    if (!auth.user()?.roles) return tabs;

    if (hasRole(ROLES.TEACHER)) {
      tabs.push("Мои ученики");
    }
    if (hasRole(ROLES.PARENT)) {
      tabs.push("Мои дети");
      tabs.push("Классы детей");
    }
    if (hasRole(ROLES.STUDENT)) {
      tabs.push("Мои родители");
      tabs.push("Мой класс");
    }
    return tabs;
  });
  const [ownerTabsActive, setOwnerTabsActive] = createSignal(ownerTabs()[0]);
  const [statusTabsActive, setStatusTabsActive] = createSignal(statusTabs[0]);

  const postsToShow = createMemo(() => {
    let posts = allPosts();

    // Status filter
    if (statusTabsActive() === statusTabs[0]) {
      posts = posts.filter((post) => !post.thingReturnedToOwner);
    } else {
      posts = posts.filter((post) => post.thingReturnedToOwner);
    }

    // Owner filter
    if (ownerTabsActive() == ownerTabs()[0]) {
      return posts;
    } else if (ownerTabsActive() == ownerTabs()[1]) {
      return posts.filter((post) => post.author.id === auth.user()!.id);
    } else if (ownerTabsActive() == ownerTabs()[2]) {
      return posts.filter((post) =>
        (studentGroupsWhereAdvisor() || [])
          .filter((group) => group?.students)
          .flatMap((group) => group.students.map((student) => student.userId))
          .flat()
          .includes(post.author.id),
      );
    } else if (ownerTabsActive() == ownerTabs()[3]) {
      return posts.filter((post) =>
        (children() || [])
          .map((student) => student.userId)
          .flat()
          .includes(post.author.id),
      );
    } else if (ownerTabsActive() == ownerTabs()[4]) {
      return posts.filter((post) =>
        (childrenStudentGroups() || [])
          .map((group) => group.students.map((student) => student.userId))
          .flat()
          .includes(post.author.id),
      );
    } else if (ownerTabsActive() == ownerTabs()[5]) {
      return posts.filter((post) =>
        (studentParents() || [])
          .map((parent) => parent.userId)
          .includes(post.author.id),
      );
    } else if (ownerTabsActive() == ownerTabs()[6]) {
      return posts.filter((post) =>
        (classmates() || [])
          .map((student) => student.userId)
          .includes(post.author.id),
      );
    }
    return posts;
  });

  return (
    <div class="max-w-4xl mx-auto space-y-6">
      <h1 class="text-2xl font-bold text-center">Объявления</h1>

      <div class="flex flex-col gap-3">
        <TabsToggle
          tabs={ownerTabs()}
          setter={setOwnerTabsActive}
          afterChange={resetAndLoad}
          tabsHTMLElementId="owner_tabs_toggle"
        />
        <TabsToggle
          tabs={statusTabs}
          setter={setStatusTabsActive}
          afterChange={resetAndLoad}
          tabsHTMLElementId="status_tabs_toggle"
        />
      </div>

      <Show when={loading() && allPosts().length === 0}>
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

          <div ref={observerRef} class="h-10">
            <Show when={postsToShow().length === 0}>
              <div class="text-center text-gray-500 py-8">
                Пока нет объявлений
              </div>
            </Show>
            <Show when={!hasMore() && postsToShow().length > 0}>
              <div class="text-center text-gray-500 py-8">
                Больше нет объявлений
              </div>
            </Show>
          </div>
        </div>
      </Show>
    </div>
  );
};

export default PublicPosts;
