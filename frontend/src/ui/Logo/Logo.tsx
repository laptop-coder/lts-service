import { JSX } from "solid-js";

import { A } from "@solidjs/router";

import { HOME_ROUTE } from "../../utils/consts";
import styles from "./logo.module.css";

const Logo = (): JSX.Element => {
  return (
    <A href={HOME_ROUTE}>
      <svg xmlns="http://www.w3.org/2000/svg" width="40" height="40">
        <image class={styles.logo} href="/src/assets/logo.svg"></image>
      </svg>
    </A>
  );
};

export default Logo;
