import { JSX, createSignal, onMount, For, createEffect } from "solid-js";
import type { Signal, Setter } from "solid-js";

import { Motion } from "solid-motionone";

const TabsToggle = (props: {
  tabs: string[];
  setter: Setter<any>;
  tabsHTMLElementId: string;
}): JSX.Element => {
  const [screenSize, setScreenSize] = createSignal({
    width: window.innerWidth,
    height: window.innerHeight,
  });

  const handleResize = () => {
    setScreenSize({
      width: window.innerWidth,
      height: window.innerHeight,
    });
  };

  const [activeTab, setActiveTab] = createSignal(0);
  const [activeTabInfo, setActiveTabInfo]: Signal<{
    left: number;
    width: number;
  }> = createSignal({ left: 0, width: 0 });
  const tabsRefs: HTMLButtonElement[] = [];

  const update = () => {
    if (screenSize()) {
      const rect = tabsRefs[activeTab()].getBoundingClientRect();
      const tabsHTMLElement = document.getElementById(props.tabsHTMLElementId);
      var tabsHTMLElementLeft = 0;
      if (tabsHTMLElement != null) {
        tabsHTMLElementLeft = tabsHTMLElement.getBoundingClientRect().left;
      }
      setActiveTabInfo({
        left: rect.left - tabsHTMLElementLeft,
        width: rect.width,
      });
    }
  };

  createEffect(() => {
    screenSize();
    update();
  });

  onMount(() => {
    setTimeout(update, 100);
    window.addEventListener("resize", handleResize);
    return () => {
      window.removeEventListener("resize", handleResize);
    };
  });

  props.setter(props.tabs[0]);
  return (
    <div class="relative bg-gray-100 h-[80px] w-full overflow-x-auto rounded-lg">
      <div
        class="flex items-center h-[80px] justify-evenly rounded-lg"
        id={props.tabsHTMLElementId}
      >
        <For each={props.tabs}>
          {(tab, index) => (
            <button
              ref={(el) => (tabsRefs[index()] = el)}
              onclick={() => {
                setActiveTab(index());
                props.setter(props.tabs[index()]);
              }}
              class="border-none bg-none p-[20px] h-full text-sm cursor-pointer flex items-center rounded-lg"
            >
              <span class="relative z-2 select-none">{tab}</span>
            </button>
          )}
        </For>
      </div>
      <Motion
        initial={false}
        animate={{
          x: activeTabInfo().left + 5,
          width: String(activeTabInfo().width - 10 + "px"),
        }}
        transition={{ duration: 0.3 }}
        class="top-[5px] bottom-[5px] rounded-lg absolute bg-white z-1"
      />
    </div>
  );
};

export default TabsToggle;
