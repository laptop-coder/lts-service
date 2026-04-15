import { JSX, createSignal, onMount, For, createEffect } from "solid-js";
import type { Signal, Setter } from "solid-js";

import { Motion } from "solid-motionone";

const TabsToggle = (props: {
  tabs: string[];
  setActiveTab: Setter<any>;
  tabsHTMLElementId: string;
  afterChange: () => void;
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

  onMount(() => {
    window.addEventListener("resize", handleResize);
    createEffect(() => {
      if (screenSize()) {
        const rect = tabsRefs[activeTab()].getBoundingClientRect();
        const tabsHTMLElement = document.getElementById(
          props.tabsHTMLElementId,
        );
        var tabsHTMLElementLeft = 0;
        if (tabsHTMLElement != null) {
          tabsHTMLElementLeft = tabsHTMLElement.getBoundingClientRect().left;
        }
        setActiveTabInfo({
          left: rect.left - tabsHTMLElementLeft,
          width: rect.width,
        });
      }
    });

    return () => {
      window.removeEventListener("resize", handleResize);
    };
  });

  props.setActiveTab(props.tabs[0]);
  return (
    <div class="relative bg-gray-100 rounded-xl w-full">
      <div
        class="relative z-10 flex justify-evenly"
        id={props.tabsHTMLElementId}
      >
        <For each={props.tabs}>
          {(tab, index) => (
            <button
              ref={(el) => (tabsRefs[index()] = el)}
              onClick={() => {
                setActiveTab(index());
                props.setActiveTab(props.tabs[index()]);
                props.afterChange();
              }}
              class="flex-1 py-4 text-sm font-medium transition-colors cursor-pointer text-center"
            >
              {tab}
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
        class="absolute top-1 bottom-1 z-0 bg-white rounded-lg shadow-sm"
      />
    </div>
  );
};

export default TabsToggle;
