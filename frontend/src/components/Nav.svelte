<script>
  import { getContext } from 'svelte';
  import { SignIn } from '../../wailsjs/go/main/App.js'

  let { appInfo } = getContext('app');
  var response;
  let { userInfo } = getContext('user');

  const signIn = async () => {
    response = await SignIn();
    $appInfo = { CanAccess: response.CanAccess, CurrentVersion: response.CurrentVersion };
    $userInfo = response.UserInfo;
  };
</script>

<nav class="bg-slate-800 mb-5 text-slate-50">
  <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
    <div class="flex h-16 items-center justify-end">
      {#if $userInfo}
        {$userInfo.Title}
      {:else}
        <button on:click={signIn}>Sign in</button>
      {/if}
    </div>
  </div>
</nav>
