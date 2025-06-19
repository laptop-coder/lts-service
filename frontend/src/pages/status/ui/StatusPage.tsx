import type { Component } from 'solid-js';
import { useParams } from '@solidjs/router';
import '../../../app/styles.css';

import { Header } from '../../../shared/ui/index';

export const StatusPage: Component = () => {
  const params = useParams();
  console.log(params.type, params.id);
  return (
    <div class='page'>
      <Header></Header>
    </div>
  );
};
