<script>
  import { navigate } from "svelte-routing";
  import {
    connections,
    activeConnectionId,
    isLsConnectionsFetched,
  } from "../stores/connections";
  import { onMount } from "svelte";
  import { activeTab, tabs } from "../stores/tabs";
  import TrashIcon from "../icons/TrashIcon.svelte";
  import { envApiEndpoint } from "../stores/env";
  import XmarkIcon from "../icons/XmarkIcon.svelte";

  let isShowAddConnectionModal = false;
  let createConnectionRequest = {
    is_existing_connection: false,
    connection_idx: 0,
    connection_name: "",
    host: "",
    port: "",
    user: "",
    password: "",
    database_name: "",
  };
  let connectionModalTitle = "New Connection";

  function openNewConnectionModal() {
    connectionModalTitle = "New Connection";
    isShowAddConnectionModal = true;
    createConnectionRequest = {
      is_existing_connection: false,
      connection_idx: 0,
      connection_name: "",
      host: "",
      port: "",
      user: "",
      password: "",
      database_name: "",
    };
  }

  function openExistingConnectionModal(connectionIdx) {
    if ($connections[connectionIdx].connection_id) {
      if ($activeConnectionId !== $connections[connectionIdx].connection_id) {
        $tabs = [];
        $activeTab = "";
      }
      $activeConnectionId = $connections[connectionIdx].connection_id;
      navigate("/table");
    }

    connectionModalTitle = "Existing Connection";
    isShowAddConnectionModal = true;
    createConnectionRequest = {
      is_existing_connection: true,
      connection_idx: connectionIdx,
      connection_name: $connections[connectionIdx].connection_name,
      host: $connections[connectionIdx].host,
      port: $connections[connectionIdx].port,
      user: $connections[connectionIdx].user,
      password: $connections[connectionIdx].password,
      database_name: $connections[connectionIdx].database_name,
    };
  }

  function closeConnectionModal() {
    isShowAddConnectionModal = false;
  }

  function saveConnection(connection) {
    let connectionIdx = connection.connection_idx;
    $connections[connectionIdx] = {
      connection_name: connection.connection_name,
      host: connection.host,
      port: connection.port,
      database_name: connection.database_name,
      user: connection.user,
      password: connection.password,
    };
    syncConnectionsToLs();
  }

  async function doCreateConnection(connection) {
    const res = await fetch(`${envApiEndpoint}/connections`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        host: connection.host,
        port: connection.port,
        database_name: connection.database_name,
        user: connection.user,
        password: connection.password,
      }),
    });
    return res;
  }

  async function connectToExistingConnection(connection) {
    let connectionIdx = connection.connection_idx;
    saveConnection(connection);

    const res = await doCreateConnection(connection);
    const jsonResp = await res.json();

    if (res.ok) {
      let connectionId = jsonResp.data.connection_id;
      $connections[connectionIdx].connection_id = connectionId;
      $activeConnectionId = connectionId;
      $tabs = [];
      $activeTab = "";
      navigate("/table");
    } else {
      return [];
    }
  }

  async function createConnection(connection) {
    if (connection.is_existing_connection) {
      connectToExistingConnection(connection);
      return;
    }

    const res = await doCreateConnection(connection);
    const jsonResp = await res.json();

    if (res.ok) {
      let connectionId = jsonResp.data.connection_id;
      $connections.push({
        connection_id: connectionId,
        connection_name: connection.connection_name,
        host: connection.host,
        port: connection.port,
        database_name: connection.database_name,
        user: connection.user,
        password: connection.password,
      });
      $connections = $connections;
      syncConnectionsToLs();
      $activeConnectionId = connectionId;
      $tabs = [];
      $activeTab = "";
      navigate("/table");
    } else {
      return [];
    }
  }

  function deleteConnection(connectionIdx) {
    $connections.splice(connectionIdx, 1);
    $connections = $connections;
    syncConnectionsToLs();
  }

  function syncConnectionsToLs() {
    let lsConnections = [];
    $connections.forEach((connection) => {
      console.log(connection);
      lsConnections.push({
        connection_name: connection.connection_name,
        host: connection.host,
        port: connection.port,
        database_name: connection.database_name,
        user: connection.user,
        password: connection.password,
      });
    });
    localStorage.setItem("connections", JSON.stringify(lsConnections));
  }

  onMount(() => {
    if (!$isLsConnectionsFetched) {
      let lsConnections = JSON.parse(localStorage.getItem("connections"));
      if (lsConnections) {
        $connections = lsConnections;
      }
      $isLsConnectionsFetched = true;
    }
  });
</script>

<div class="container">
  <div class="connection-header">
    <h1 class="connection-title">Connections</h1>
    <button
      class="connection-btn-add"
      on:click={() => openNewConnectionModal()}
    >
      Add Connection
    </button>
  </div>
  <hr />
  <div class="connection-modal {isShowAddConnectionModal ? 'active' : ''}">
    <div class="connection-form-container">
      <button
        class="connection-modal-btn-close"
        on:click={() => closeConnectionModal()}
      >
        <XmarkIcon />
      </button>
      <h3 class="connection-form-title">{connectionModalTitle}</h3>
      <form
        on:submit|preventDefault={() =>
          createConnection(createConnectionRequest)}
        class="connection-form"
      >
        <div class="connection-form-item">
          <label for="connection_name" class="connection-form-label">
            Connection Name
          </label>
          <input
            class="connection-form-input"
            type="text"
            name="connection_name"
            id="connection_name"
            bind:value={createConnectionRequest.connection_name}
          />
        </div>
        <div class="connection-form-item">
          <label for="host" class="connection-form-label">Host</label>
          <input
            class="connection-form-input"
            type="text"
            name="host"
            id="host"
            bind:value={createConnectionRequest.host}
          />
        </div>
        <div class="connection-form-item">
          <label for="port" class="connection-form-label">Port</label>
          <input
            class="connection-form-input"
            type="text"
            name="port"
            id="port"
            bind:value={createConnectionRequest.port}
          />
        </div>
        <div class="connection-form-item">
          <label for="user" class="connection-form-label">User</label>
          <input
            class="connection-form-input"
            type="text"
            name="user"
            id="user"
            bind:value={createConnectionRequest.user}
          />
        </div>
        <div class="connection-form-item">
          <label for="password" class="connection-form-label">Password</label>
          <input
            class="connection-form-input"
            type="password"
            name="password"
            id="password"
            bind:value={createConnectionRequest.password}
          />
        </div>
        <div class="connection-form-item">
          <label for="database_name" class="connection-form-label">
            Database Name
          </label>
          <input
            class="connection-form-input"
            type="text"
            name="database_name"
            id="database_name"
            bind:value={createConnectionRequest.database_name}
          />
        </div>
        <div class="connection-form-footer">
          <button
            class="connection-btn-save"
            on:click|preventDefault={() =>
              saveConnection(createConnectionRequest)}
          >
            Save
          </button>
          <button type="submit" class="connection-btn-connect">Connect</button>
        </div>
      </form>
    </div>
  </div>

  <div class="connection-list">
    {#if $connections.length != 0}
      {#each $connections as connection, connectionIdx}
        <div class="connection-item">
          <button
            on:click={() => openExistingConnectionModal(connectionIdx)}
            class="connection-open-btn"
          >
            {connection.connection_name}
          </button>
          <div class="connection-item-info">
            {#if connection.connection_id > 0}
              <div class="connection-active-sign" />
            {/if}
            <button
              on:click={() => deleteConnection(connectionIdx)}
              class="connection-del-btn"
            >
              <TrashIcon height="0.7rem" fill="#555" />
            </button>
          </div>
        </div>
      {/each}
    {/if}
  </div>
</div>

<style>
  .container {
    padding: 1rem 2rem;
    font-size: 0.9rem;
  }

  .connection-modal {
    z-index: 1000;
    display: none;
    position: absolute;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    justify-content: center;
    background-color: rgba(0, 0, 0, 0.7);
  }

  .connection-modal.active {
    display: flex;
  }

  .connection-modal-btn-close {
    position: absolute;
    top: 1rem;
    right: 1rem;
    padding: 0.4rem 0.5rem;
    border-radius: 50%;
    border: 1px solid #777;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .connection-form-container {
    position: relative;
    margin-top: 4rem;
    width: 500px;
    max-height: 600px;
    overflow: scroll;
    background-color: white;
    padding: 1rem 1.5rem;
    border-radius: 1rem;
  }

  .connection-form {
    margin-top: 0.5rem;
  }

  .connection-form-title {
    font-size: 1.5rem;
  }

  .connection-form-item {
    display: flex;
    gap: 0.5rem;
    align-items: center;
    padding: 0.4rem 0;
  }

  .connection-form-label {
    display: block;
    min-width: 140px;
  }

  .connection-form-input {
    padding: 0.3rem 0.5rem;
    border-radius: 0.4rem;
    border: 1px solid #999;
    display: block;
    width: 100%;
  }

  .connection-form-footer {
    margin-top: 0.7rem;
    display: flex;
    gap: 0.4rem;
    justify-content: end;
  }

  .connection-header {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding-bottom: 0.3rem;
  }

  .connection-title {
    font-size: 2rem;
    font-weight: 600;
  }

  .connection-btn-add,
  .connection-btn-connect,
  .connection-btn-save {
    padding: 0.3rem 0.5rem;
  }

  .connection-list {
    display: flex;
    gap: 0.7rem;
    flex-wrap: wrap;
  }

  .connection-item {
    position: relative;
  }

  .connection-item-info {
    position: absolute;
    top: 0.4rem;
    right: 0.4rem;
    display: flex;
    align-items: center;
    gap: 0.4rem;
  }

  .connection-active-sign {
    width: 9px;
    height: 9px;
    border-radius: 50%;
    background-color: green;
  }

  .connection-open-btn {
    padding: 1rem 4rem 1rem 1rem;
  }

  .connection-del-btn {
    width: 20px;
    height: 20px;
    border-radius: 0.3rem;
    border: 1px solid #777;
  }

  hr {
    border-bottom: 1px solid #555;
    margin-bottom: 2rem;
  }
</style>
