import { JSX, createSignal, onMount, For, createEffect } from 'solid-js';
import type { Signal, Setter } from 'solid-js';

import TabsToggleWrapper from '../TabsToggleWrapper/TabsToggleWrapper';
import styles from './TabsToggle.module.css';

import { Motion } from 'solid-motionone';

const TabsToggle = (props: {
  tabs: string[];
  tabsNames: string[];
  setter: Setter<any>;
  tabsHTMLElementId: string;
  fullscreen?: boolean;
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
    window.addEventListener('resize', handleResize);
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
      window.removeEventListener('resize', handleResize);
    };
  });

  props.setter(props.tabs[0]);
  return (
    <TabsToggleWrapper fullscreen={props.fullscreen}>
      <div
        class={styles.tabs}
        id={props.tabsHTMLElementId}
      >
        <For each={props.tabs}>
          {(_, index) => (
            <button
              ref={(el) => (tabsRefs[index()] = el)}
              onclick={() => {
                setActiveTab(index());
                props.setter(props.tabs[index()]);
              }}
              class={styles.tab}
            >
              <span class={styles.tab_text}>{props.tabsNames[index()]}</span>
            </button>
          )}
        </For>
      </div>
      <Motion
        initial={false}
        animate={{
          x: activeTabInfo().left + 5,
          width: String(activeTabInfo().width - 10 + 'px'),
        }}
        transition={{ duration: 0.3 }}
        class={styles.pointer}
      />
    </TabsToggleWrapper>
  );
};

export default TabsToggle;
