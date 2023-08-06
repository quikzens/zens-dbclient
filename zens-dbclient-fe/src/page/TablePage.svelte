<script>
  import { navigate } from "svelte-routing";
  import { tabs, activeTab } from "../stores/tabs";
  import { activeConnectionId, connections } from "../stores/connections";
  import { onMount } from "svelte";
  import WifiIcon from "../icons/WifiIcon.svelte";
  import { envApiEndpoint } from "../stores/env";
  import PowerOffIcon from "../icons/PowerOffIcon.svelte";

  onMount(() => {
    if ($activeConnectionId === 0) {
      navigate("/connection");
    }
  });

  async function getTables() {
    if ($activeConnectionId !== 0) {
      const res = await fetch(
        `${envApiEndpoint}/${$activeConnectionId}/tables`
      );
      const jsonResp = await res.json();

      if (res.ok) {
        return jsonResp.data;
      } else {
        return [];
      }
    } else {
      return [];
    }
  }

  async function closeActiveConnection() {
    if ($activeConnectionId !== 0) {
      const res = await fetch(
        `${envApiEndpoint}/connections/${$activeConnectionId}`,
        {
          method: "DELETE",
        }
      );
      const jsonResp = await res.json();

      if (res.ok) {
        let connectionIdx = $connections.findIndex(
          (c) => c.connection_id === $activeConnectionId
        );
        $activeConnectionId = 0;
        $connections[connectionIdx].connection_id = 0;
        navigate("/connection");
        return jsonResp.data;
      } else {
        return [];
      }
    } else {
      return [];
    }
  }

  function addRecordTabIfNotExist(tabName) {
    if (!$tabs.includes(tabName)) {
      $tabs.push(tabName);
    }
  }

  function visitRecordTab(tabName) {
    addRecordTabIfNotExist(tabName);
    $activeTab = tabName;
    navigate("/record");
  }
</script>

<div class="container">
  {#await getTables() then tables}
    <div class="table-list-box">
      <h1>Open a Table</h1>
      <div class="table-list">
        {#each tables as table}
          <button class="table-item" on:click={() => visitRecordTab(table)}>
            {table}
          </button>
        {/each}
      </div>
    </div>
  {/await}
</div>
<div class="connection-btn">
  <button class="connection-btn-close" on:click={() => closeActiveConnection()}>
    <PowerOffIcon height="0.9rem" fill="#222" />
  </button>
  <button class="connection-btn-list" on:click={() => navigate("/connection")}>
    <WifiIcon height="0.75rem" fill="#222" />
  </button>
</div>

<style>
  .container {
    display: flex;
    justify-content: center;
  }

  .table-list-box {
    margin: 3rem;
    padding: 1rem;
    width: 500px;
    border: 1px solid #777;
    border-radius: 0.5rem;
  }

  .table-list-box h1 {
    font-size: 1.2rem;
  }

  .table-list {
    padding-top: 1rem;
  }

  .table-item {
    text-align: left;
    width: 100%;
    font-size: 0.9rem;
    display: block;
    padding: 0.3rem 0.5rem;
    transition: 0.3s;
    border: none;
    border-radius: 0.3rem;
    background-color: transparent;
  }

  .table-item:hover {
    background-color: #d9d9d9;
  }

  .connection-btn {
    position: absolute;
    top: 1rem;
    right: 1rem;
    display: flex;
    gap: 0.3rem;
  }

  .connection-btn-list,
  .connection-btn-close {
    display: flex;
    justify-content: center;
    align-items: center;
    padding: 0.3rem;
  }
</style>
