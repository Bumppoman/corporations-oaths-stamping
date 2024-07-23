<script>
  import { getContext, setContext } from 'svelte';
  import { writable } from 'svelte/store';

  import FileList from './components/FileList.svelte';
  import Nav from './components/Nav.svelte';

  setContext('app', {
    appInfo: writable(null),
  });

  setContext('user', {
    userInfo: writable(null),
  });

  const { appInfo } = getContext('app');
  const { userInfo } = getContext('user');
</script>

<main>
  <Nav />

  <div class="mx-auto px-5">
    {#if $userInfo && $appInfo && $appInfo.CanAccess}
      {#if $appInfo && $appInfo.CurrentVersion == '1.0.0'}
        <FileList />
      {:else}
        <div class="flex items-center justify-center min-h-screen text-slate-50 w-full">
          <div>Your application is out of date.  Please download the current version.</div>
        </div>
      {/if}
    {:else if $appInfo && !$appInfo.CanAccess}
      <div class="flex items-center justify-center min-h-screen text-slate-50 w-full">
        <div>You do not have access to this application.  Contact <a class="underline" href="mailto:dos.sm.data@dos.ny.gov">data management</a>.</div>
      </div>
    {:else}
      <div class="flex items-center justify-center min-h-screen text-slate-50 w-full">
        <div>Please sign in to use this application.</div>
      </div>
    {/if}
  </div>
</main>

