import { JSX } from 'solid-js';

import Header from '../components/Header/Header';
import Content from '../components/Content/Content';
import Footer from '../components/Footer/Footer';
import Page from '../ui/Page/Page';
import SquareImageButton from '../ui/SquareImageButton/SquareImageButton';

const HomePage = (): JSX.Element => {
  return (
    <Page>
      <Header>
        <SquareImageButton>
          <img src='/src/assets/add.svg'></img>
        </SquareImageButton>
        <SquareImageButton>
          <img src='/src/assets/reload.svg'></img>
        </SquareImageButton>
      </Header>
      <Content></Content>
      <Footer></Footer>
    </Page>
  );
};

export default HomePage;
