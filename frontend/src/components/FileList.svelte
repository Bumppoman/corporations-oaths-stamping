<script>
  import { onMount } from 'svelte';
  import { DownloadAttachment, LoadUnstamped, UploadStamped } from '../../wailsjs/go/main/App.js'
  import { stampPDF } from '../stamping.js'

  export let lastUpdated = new Date().toLocaleString();
  let unstamped = [];

  $: refresh = async () => {
    const rawUnstamped = await LoadUnstamped() || [];
    unstamped  = rawUnstamped.map(pdf => ({ ...pdf, Selected: true }));
    lastUpdated = new Date().toLocaleString();
  };

  onMount(() => refresh());

  const stamp = async () => {
    for (const pdf of unstamped.filter(pdf => pdf.Selected)) {
      updatePDFStatus(pdf, 'Downloading');
      let blob = await DownloadAttachment(pdf.Id);

      // Retry download once
      if (!blob) {
        blob = await DownloadAttachment(pdf.Id);
        if (!blob) {
          updatePDFStatus(pdf, 'Error');
          continue;
        }
      }

      // Stamp PDF (unrecoverable on failure)
      updatePDFStatus(pdf, 'Stamping');
      const stamped = await stampPDF(Uint8Array.from(atob(blob), c => c.charCodeAt(0)), pdf.StampText);
      if (!stamped) {
        updatePDFStatus(pdf, 'Error');
        continue;
      }

      // Upload stamped PDF
      updatePDFStatus(pdf, 'Uploading');
      const base64PDF = [];
      for (const byte of stamped) {
        base64PDF.push(String.fromCharCode(byte));
      }
      const error = await UploadStamped(pdf.Id, btoa(base64PDF.join('')));

      if (error) {
        updatePDFStatus(pdf, 'Error');
        continue;
      }

      updatePDFStatus(pdf, 'Complete');
    }

    refresh();
  };

  const updatePDFStatus = (pdf, status) => {
    pdf.status = status;

    // For reactivity purposes
    unstamped = unstamped;
  };
</script>

<div class="bg-gray-900">
  <div class="mx-auto max-w-7xl">
    <div class="bg-gray-900 py-10">
      <div class="px-4 sm:px-6 lg:px-8">
        <div class="sm:flex sm:items-center">
          <div class="sm:flex-auto">
            <h1 class="text-base font-semibold leading-6 text-white">Unstamped PDFs</h1>
          </div>
          <div class="flex gap-x-3 mt-4">
            <button
              type="button"
              class="block rounded-md bg-white px-3 py-2 text-center text-sm font-semibold ring-1 ring-inset ring-slate-300 shadow-sm text-slate-900 hover:bg-slate-50"
              on:click={refresh}
            >
              Refresh
            </button>
            <button
              type="button"
              class="block rounded-md bg-sky-500 px-3 py-2 text-center text-sm font-semibold text-white hover:bg-sky-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-sky-500"
              on:click={stamp}
            >
              Stamp PDFs
            </button>
          </div>
        </div>
        <div class="mt-8 flow-root">
          <div class="-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
            <div class="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
              <table class="min-w-full divide-y divide-gray-700">
                <thead>
                  <tr>
                    <th scope="col"></th>
                    <th
                      scope="col"
                      class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-white sm:pl-0"
                    >
                      Submitter
                    </th>
                    <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-white">
                      Submission Date
                    </th>
                    <th
                      scope="col"
                      class="px-3 py-3.5 text-left text-sm font-semibold text-white w-1/3"
                    >
                      Status
                    </th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-slate-800">
                  <tr class="hidden only:table-row">
                    <td colspan="5" class="py-4 pl-4 text-sm font-medium text-white sm:pl-0">
                      There are currently no PDFs to stamp.
                    </td>
                  </tr>
                  {#each unstamped as pdf}
                    <tr>
                      <td>
                        <input
                          type="checkbox"
                          class="h-3 w-3 rounded border-slate-300 text-sky-600 focus:ring-sky-600"
                          bind:checked={pdf.Selected}
                        />
                      </td>
                      <td
                        class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-white sm:pl-0"
                      >
                        {pdf.Title}
                      </td>
                      <td class="whitespace-nowrap px-3 py-4 text-sm text-slate-300">
                        {new Date(pdf.Created).toLocaleDateString()}
                      </td>
                      <td class="whitespace-nowrap px-3 py-4 text-sm text-slate-300">
                        {pdf.status || 'Pending'}
                      </td>
                    </tr>
                  {/each}
                </tbody>
                <tfoot>
                  <tr>
                    <td colspan="4" class="py-4 pl-4 text-xs font-medium text-right text-white sm:pl-0">
                      Last updated {lastUpdated}
                    </td>
                  </tr>
                </tfoot>
              </table>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>
