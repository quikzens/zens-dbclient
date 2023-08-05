<script>
  import { tabs, activeTab } from "../stores/tabs.js";
  import { navigate } from "svelte-routing";
  import { activeConnectionId, connections } from "../stores/connections";
  import { envApiEndpoint } from "../stores/env";
  import RecordTab from "./RecordTab.svelte";
  import PlusIcon from "../icons/PlusIcon.svelte";
  import XmarkIcon from "../icons/XmarkIcon.svelte";
  import WifiIcon from "../icons/WifiIcon.svelte";
  import PowerOffIcon from "../icons/PowerOffIcon.svelte";

  function goToTab(tabName) {
    $activeTab = tabName;
  }

  function deleteRecordTab(tabName) {
    let deletedTabIdx = $tabs.findIndex((tab) => tab === tabName);
    $tabs.splice(deletedTabIdx, 1);
    if ($tabs.length === 0) {
      navigate("/table");
    }
    $tabs = $tabs;
    $activeTab = $tabs[0];
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
</script>

<div class="tabs">
  <div class="tab-links">
    <div>
      {#each $tabs as tab}
        <div class="tab-link">
          <button
            class="tab-link-visit-btn {tab === $activeTab ? 'active' : ''}"
            on:click={() => goToTab(tab)}
          >
            {tab}
          </button>
          <button
            class="tab-link-close-btn"
            on:click={() => deleteRecordTab(tab)}
          >
            <XmarkIcon height="0.75rem" fill="#444" />
          </button>
        </div>
      {/each}
      <button class="add-tab-btn" on:click={() => navigate("/table")}>
        <PlusIcon height="1rem" fill="#444" />
      </button>
    </div>
    <div>
      <button
        class="connection-btn-close"
        on:click={() => closeActiveConnection()}
      >
        <PowerOffIcon height="0.9rem" fill="#222" />
      </button>
      <button
        class="connection-btn-list"
        on:click={() => navigate("/connection")}
      >
        <WifiIcon height="0.75rem" fill="#222" />
      </button>
    </div>
  </div>
  <div class="tab-panels">
    {#each $tabs as tab}
      <div class="tab-panel {tab === $activeTab ? 'active' : ''}">
        <RecordTab table_name={tab} />
      </div>
    {/each}
  </div>
</div>

<style>
  .tabs {
    font-size: 0.9rem;
  }

  .tab-links {
    display: flex;
    justify-content: space-between;
    padding: 0.5rem 1rem;
    border-bottom: 1px solid #cdcdcd;
  }

  .tab-links > div {
    display: flex;
    gap: 0.4rem;
  }

  .tab-panel {
    display: none;
  }

  .tab-panel.active {
    display: block;
  }

  .tab-link {
    position: relative;
    display: flex;
  }

  .tab-link-visit-btn {
    padding: 0.25rem 1.5rem 0.25rem 0.75rem;
    background-color: transparent;
    transition: 0.3s;
    border: 1px solid #888;
    border-radius: 0.5rem;
  }

  .tab-link-visit-btn:hover,
  .tab-link-visit-btn.active {
    background-color: #e4e4e4;
  }

  .tab-link-close-btn {
    width: 15px;
    height: 15px;
    padding: 0.05rem;
    border-radius: 50%;
    border: none;
    position: absolute;
    top: 50%;
    right: 5px;
    display: flex;
    justify-content: center;
    align-items: center;
    transform: translateY(-50%);
  }

  .add-tab-btn {
    display: flex;
    justify-content: center;
    align-items: center;
    padding: 0.3rem 0.35rem;
  }

  .connection-btn-list,
  .connection-btn-close {
    display: flex;
    justify-content: center;
    align-items: center;
    padding: 0.3rem;
  }
</style>
