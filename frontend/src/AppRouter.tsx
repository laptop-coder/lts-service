import { JSX } from 'solid-js';

import { Router, Route } from '@solidjs/router';

import HomePage from './pages/HomePage';
import AddThingPage from './pages/AddThingPage';
import ThingStatusPage from './pages/ThingStatusPage';
import NotFoundPage from './pages/errors/NotFoundPage/NotFoundPage';
import ModeratorRegisterPage from './pages/ModeratorRegisterPage';
import ModeratorLoginPage from './pages/ModeratorLoginPage';
import {
  HOME_ROUTE,
  ADD_THING_ROUTE,
  THING_STATUS_ROUTE,
  MODERATOR_REGISTER_ROUTE,
  MODERATOR_LOGIN_ROUTE,
} from './utils/consts';

const AppRouter = (): JSX.Element => {
  return (
    <Router>
      <Route
        path={ADD_THING_ROUTE}
        component={AddThingPage}
      />
      <Route
        path={THING_STATUS_ROUTE}
        component={ThingStatusPage}
      />
      <Route
        path={HOME_ROUTE}
        component={HomePage}
      />
      <Route
        path={MODERATOR_REGISTER_ROUTE}
        component={ModeratorRegisterPage}
      />
      <Route
        path={MODERATOR_LOGIN_ROUTE}
        component={ModeratorLoginPage}
      />
      <Route
        path='*404'
        component={NotFoundPage}
      />
    </Router>
  );
};

export default AppRouter;
