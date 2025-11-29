import { JSX } from 'solid-js';

import { Router, Route } from '@solidjs/router';

import HomePage from './pages/Home';
import LoginUserPage from './pages/LoginUser';
import RegisterUserPage from './pages/RegisterUser';
import LoginModeratorPage from './pages/LoginModerator';
import RegisterModeratorPage from './pages/RegisterModerator';
import UserProfilePage from './pages/UserProfile';
import UserThingAddPage from './pages/UserThingAdd';
import UserThingEditPage from './pages/UserThingEdit';
import UserThingStatusPage from './pages/UserThingStatus';
import ModeratorHomePage from './pages/ModeratorHome';
import ModeratorProfilePage from './pages/ModeratorProfile';

import {
  HOME__ROUTE,
  LOGIN_USER__ROUTE,
  REGISTER_USER__ROUTE,
  LOGIN_MODERATOR__ROUTE,
  REGISTER_MODERATOR__ROUTE,
  USER__PROFILE__ROUTE,
  USER__THING_ADD__ROUTE,
  USER__THING_EDIT__ROUTE,
  USER__THING_STATUS__ROUTE,
  MODERATOR__HOME__ROUTE,
  MODERATOR__PROFILE__ROUTE,
} from './utils/consts';

const AppRouter = (): JSX.Element => {
  return (
    <Router>
      {/*For all users*/}
      <Route
        path={HOME__ROUTE}
        component={HomePage}
      />
      <Route
        path={LOGIN_USER__ROUTE}
        component={LoginUserPage}
      />
      <Route
        path={LOGIN_MODERATOR__ROUTE}
        component={LoginModeratorPage}
      />
      <Route
        path={REGISTER_USER__ROUTE}
        component={RegisterUserPage}
      />
      <Route
        path={REGISTER_MODERATOR__ROUTE}
        component={RegisterModeratorPage}
      />

      {/*For registered users*/}
      <Route
        path={USER__PROFILE__ROUTE}
        component={UserProfilePage}
      />
      <Route
        path={USER__THING_ADD__ROUTE}
        component={UserThingAddPage}
      />
      <Route
        path={USER__THING_EDIT__ROUTE}
        component={UserThingEditPage}
      />
      <Route
        path={USER__THING_STATUS__ROUTE}
        component={UserThingStatusPage}
      />

      {/*For moderator*/}
      <Route
        path={MODERATOR__HOME__ROUTE}
        component={ModeratorHomePage}
      />
      <Route
        path={MODERATOR__PROFILE__ROUTE}
        component={ModeratorProfilePage}
      />
    </Router>
  );
};

export default AppRouter;
