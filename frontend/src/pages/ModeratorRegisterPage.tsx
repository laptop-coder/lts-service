import { JSX } from 'solid-js';

import ModeratorRegisterForm from '../ui/ModeratorRegisterForm';
import ModeratorAuthPage from '../ui/ModeratorAuthPage/ModeratorAuthPage';

const ModeratorRegisterPage = (): JSX.Element => {
  return (
    <ModeratorAuthPage>
      <ModeratorRegisterForm />
    </ModeratorAuthPage>
  );
};

export default ModeratorRegisterPage;
