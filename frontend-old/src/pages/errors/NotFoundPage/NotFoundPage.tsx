import { JSX } from 'solid-js';

import Content from '../../../components/Content/Content';
import Footer from '../../../components/Footer/Footer';
import Header from '../../../components/Header/Header';
import Page from '../../../ui/Page/Page';
import styles from './NotFoundPage.module.css';
import { ASSETS_ROUTE } from '../../../utils/consts';

const NotFoundPage = (): JSX.Element => {
  document.title = 'Страница не найдена';
  return (
    <Page>
      <Header />
      <Content class={styles.not_found_page_content}>
        <img
          src={`${ASSETS_ROUTE}/404.png`}
          class={styles.not_found_page_image}
        />
        <span class={styles.not_found_page_text}>
          Кажется, эта страница потерялась!
        </span>
      </Content>
      <Footer />
    </Page>
  );
};

export default NotFoundPage;
