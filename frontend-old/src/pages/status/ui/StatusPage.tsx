import '@/app/styles.css';
import type { Component } from 'solid-js';
import { Header } from '@/shared/ui/index';
import { useSearchParams } from '@solidjs/router';

export const StatusPage: Component = () => {
  const [searchParams, setSearchParams] = useSearchParams();
  console.log(searchParams.type, searchParams.id);
  return (
    <div class='page'>
      <Header></Header>
    </div>
  );
};
