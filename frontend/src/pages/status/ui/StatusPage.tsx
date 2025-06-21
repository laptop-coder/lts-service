import type { Component } from 'solid-js';
import { useSearchParams } from '@solidjs/router';
import '../../../app/styles.css';

import { Header } from '../../../shared/ui/index';

export const StatusPage: Component = () => {
  const [searchParams, setSearchParams] = useSearchParams();
  console.log(searchParams.type, searchParams.id);
  return (
    <div class='page'>
      <Header></Header>
    </div>
  );
};
