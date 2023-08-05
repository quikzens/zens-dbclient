<script>
  import ChevronLeftIcon from "../icons/ChevronLeftIcon.svelte";
  import ChevronRightIcon from "../icons/ChevronRightIcon.svelte";
  import CircleMinusIcon from "../icons/CircleMinusIcon.svelte";
  import CirclePlusIcon from "../icons/CirclePlusIcon.svelte";
  import CheckIcon from "../icons/CheckIcon.svelte";
  import SortIcon from "../icons/SortIcon.svelte";
  import { activeConnectionId } from "../stores/connections";

  export let table_name;

  const envApiEndpoint = import.meta.env.VITE_API_ENDPOINT;

  let recordPerPage = 25;
  const operators = [
    "=",
    "!=",
    "<",
    ">",
    "<=",
    ">=",
    "IN",
    "NOT IN",
    "LIKE",
    "BETWEEN",
    "IS NULL",
    "IS NOT NULL",
  ];
  const oneColumnOperators = [
    "=",
    "!=",
    "<",
    ">",
    "<=",
    ">=",
    "IN",
    "NOT IN",
    "LIKE",
  ];
  const twoColumnOperators = ["BETWEEN"];
  let filterSettings = [];
  let activeColumns = [];

  let tableName = table_name;
  let firstColumnName = "";
  let sortColumn = "";
  let sortType = "";
  let resultCount = 0;
  let limit = recordPerPage;
  let offset = 0;
  let isPrevPageDisable = true;
  let recordsPromise = getTableRecords();
  let isShowFilterSetting = false;
  let isShowColumnSetting = false;

  $: {
    // handle limit change
    limit;
    refetchTableRecords();
  }

  $: {
    // handle offset change
    offset;
    refetchTableRecords();
  }

  function refetchTableRecords(filterSettings) {
    recordsPromise = getTableRecords(filterSettings);
  }

  function prevPagination() {
    offset -= recordPerPage;
  }

  function nextPagination() {
    offset += recordPerPage;
  }

  function setSort(columnName) {
    if (sortColumn !== columnName) {
      sortType = "";
    }

    switch (sortType) {
      case "":
        sortColumn = columnName;
        sortType = "asc";
        break;
      case "asc":
        sortColumn = columnName;
        sortType = "desc";
        break;
      case "desc":
        sortColumn = "";
        sortType = "";
        break;
    }

    refetchTableRecords();
  }

  function toggleShowFilterSetting() {
    isShowFilterSetting = !isShowFilterSetting;
    if (isShowFilterSetting && filterSettings.length === 0) {
      filterSettings = [
        {
          is_active: true,
          field: firstColumnName,
          operator: "=",
          first_value: "",
          second_value: "",
        },
      ];
    }
    isShowColumnSetting = false;
  }

  function toggleShowColumnSetting() {
    isShowColumnSetting = !isShowColumnSetting;
    isShowFilterSetting = false;
  }

  function addFilterSetting() {
    filterSettings.push({
      is_active: true,
      field: firstColumnName,
      operator: "=",
      first_value: "",
      second_value: "",
    });
    filterSettings = filterSettings;
  }

  function subtractFilterSetting(index) {
    filterSettings.splice(index, 1);
    if (filterSettings.length === 0) {
      isShowFilterSetting = false;
      refetchTableRecords(filterSettings);
    }
    filterSettings = filterSettings;
  }

  function applyFilterSetting() {
    // filter active filterSettings
    let activeFilterSettings = filterSettings.filter(function (filterSetting) {
      return filterSetting.is_active;
    });
    refetchTableRecords(activeFilterSettings);
  }

  function toggleColumnVisibility(index) {
    activeColumns[index].is_hide = !activeColumns[index].is_hide;
  }

  async function getTableColumns() {
    const res = await fetch(
      `${envApiEndpoint}/${$activeConnectionId}/tables/${tableName}/columns`
    );
    const jsonResp = await res.json();

    if (res.ok) {
      firstColumnName = jsonResp.data[0].column_name;
      jsonResp.data.forEach((column) => {
        activeColumns.push({
          column_name: column.column_name,
          data_type: column.data_type,
          is_hide: false,
        });
      });
      return jsonResp.data;
    } else {
      return [];
    }
  }

  async function getTableRecords(filterSettings) {
    if (offset - recordPerPage < 0) {
      isPrevPageDisable = true;
    } else {
      isPrevPageDisable = false;
    }

    const res = await fetch(
      `${envApiEndpoint}/${$activeConnectionId}/tables/${tableName}/records?sort_by=${sortColumn}&order_by=${sortType}&limit=${limit}&offset=${offset}`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          conditions: filterSettings,
        }),
      }
    );
    const jsonResp = await res.json();

    if (res.ok) {
      resultCount = jsonResp.data.length;
      return jsonResp.data;
    } else {
      return [];
    }
  }
</script>

<div class="container">
  {#await getTableColumns() then columns}
    <div class="configuration">
      <div class="filter-section">
        <button class="filter-btn" on:click={() => toggleShowFilterSetting()}>
          <svg
            xmlns="http://www.w3.org/2000/svg"
            height="1em"
            viewBox="0 0 512 512"
          >
            <!--! Font Awesome Free 6.4.0 by @fontawesome - https://fontawesome.com License - https://fontawesome.com/license (Commercial License) Copyright 2023 Fonticons, Inc. -->
            <path
              d="M3.9 54.9C10.5 40.9 24.5 32 40 32H472c15.5 0 29.5 8.9 36.1 22.9s4.6 30.5-5.2 42.5L320 320.9V448c0 12.1-6.8 23.2-17.7 28.6s-23.8 4.3-33.5-3l-64-48c-8.1-6-12.8-15.5-12.8-25.6V320.9L9 97.3C-.7 85.4-2.8 68.8 3.9 54.9z"
            />
          </svg>
          <span>Filters</span>
        </button>
        <div class="column-toggle-wrapper">
          <button class="column-btn" on:click={() => toggleShowColumnSetting()}>
            <svg
              xmlns="http://www.w3.org/2000/svg"
              height="1em"
              viewBox="0 0 512 512"
            >
              <!--! Font Awesome Free 6.4.0 by @fontawesome - https://fontawesome.com License - https://fontawesome.com/license (Commercial License) Copyright 2023 Fonticons, Inc. -->
              <path
                d="M0 416c0 17.7 14.3 32 32 32l54.7 0c12.3 28.3 40.5 48 73.3 48s61-19.7 73.3-48L480 448c17.7 0 32-14.3 32-32s-14.3-32-32-32l-246.7 0c-12.3-28.3-40.5-48-73.3-48s-61 19.7-73.3 48L32 384c-17.7 0-32 14.3-32 32zm128 0a32 32 0 1 1 64 0 32 32 0 1 1 -64 0zM320 256a32 32 0 1 1 64 0 32 32 0 1 1 -64 0zm32-80c-32.8 0-61 19.7-73.3 48L32 224c-17.7 0-32 14.3-32 32s14.3 32 32 32l246.7 0c12.3 28.3 40.5 48 73.3 48s61-19.7 73.3-48l54.7 0c17.7 0 32-14.3 32-32s-14.3-32-32-32l-54.7 0c-12.3-28.3-40.5-48-73.3-48zM192 128a32 32 0 1 1 0-64 32 32 0 1 1 0 64zm73.3-64C253 35.7 224.8 16 192 16s-61 19.7-73.3 48L32 64C14.3 64 0 78.3 0 96s14.3 32 32 32l86.7 0c12.3 28.3 40.5 48 73.3 48s61-19.7 73.3-48L480 128c17.7 0 32-14.3 32-32s-14.3-32-32-32L265.3 64z"
              />
            </svg>
            <span>Columns</span>
          </button>
          <div class="column-toggle {isShowColumnSetting ? 'is-active' : ''}">
            <div class="column-toggle-list">
              {#each activeColumns as column, index}
                <button
                  class="column-toggle-item"
                  on:click={() => toggleColumnVisibility(index)}
                >
                  <span class="column-toggle-icon">
                    {#if !column.is_hide}
                      <CheckIcon />
                    {/if}
                  </span>
                  <span>
                    {column.column_name}
                  </span>
                </button>
              {/each}
            </div>
          </div>
        </div>
      </div>
      <div class="pagination-section">
        <p>{resultCount} rows</p>
        <div class="pagination">
          <button
            class="pagination-prev-btn"
            on:click={() => prevPagination()}
            disabled={isPrevPageDisable}
          >
            <ChevronLeftIcon />
          </button>
          <div class="pagination-input">
            <input
              type="number"
              name="limit"
              id="pagination-input-limit"
              bind:value={limit}
            />
            <input
              type="number"
              name="offset"
              id="pagination-input-offset"
              bind:value={offset}
            />
          </div>
          <button class="pagination-next-btn" on:click={() => nextPagination()}>
            <ChevronRightIcon />
          </button>
        </div>
      </div>
    </div>

    {#if isShowFilterSetting}
      <div class="filter-setting">
        <div class="filter-setting-items">
          {#each filterSettings as filterSetting, index}
            <div class="filter-setting-item">
              <div class="filter-setting-item-left">
                <input
                  type="checkbox"
                  name="active"
                  class="filter-item-active"
                  bind:checked={filterSetting.is_active}
                  on:change={() => applyFilterSetting()}
                />
                <div class="filter-item-select-wrapper">
                  <select
                    name="field"
                    class="filter-item-field"
                    bind:value={filterSetting.field}
                  >
                    {#each columns as column}
                      <option value={column.column_name}
                        >{column.column_name}</option
                      >
                    {/each}
                  </select>
                  <div class="filter-item-select-icon">
                    <SortIcon />
                  </div>
                </div>
                <div class="filter-item-select-wrapper">
                  <select
                    name="operator"
                    class="filter-item-operator"
                    bind:value={filterSetting.operator}
                  >
                    {#each operators as operator}
                      <option value={operator}>{operator}</option>
                    {/each}
                  </select>
                  <div class="filter-item-select-icon">
                    <SortIcon />
                  </div>
                </div>
                {#if oneColumnOperators.includes(filterSetting.operator) || twoColumnOperators.includes(filterSetting.operator)}
                  <input
                    type="text"
                    name="value"
                    class="filter-item-value"
                    bind:value={filterSetting.first_value}
                  />
                {/if}
                {#if twoColumnOperators.includes(filterSetting.operator)}
                  <input
                    type="text"
                    name="value"
                    class="filter-item-value"
                    bind:value={filterSetting.second_value}
                  />
                {/if}
              </div>
              <button
                class="filter-item-btn"
                on:click={() => subtractFilterSetting(index)}
              >
                <CircleMinusIcon />
              </button>
            </div>
          {/each}
        </div>
        <div class="filter-setting-submit">
          <button class="filter-item-btn" on:click={() => addFilterSetting()}>
            <CirclePlusIcon />
          </button>
          <button
            type="submit"
            on:click={() => applyFilterSetting()}
            class="filter-setting-submit-btn">Apply Filter</button
          >
        </div>
      </div>
    {/if}

    <div class="table-container">
      <div>
        <table class="table">
          <tr class="table-header">
            {#each columns as column, index}
              {#if !activeColumns[index].is_hide}
                <th class="table-cell">
                  <div
                    class="table-cell-btn"
                    on:click={() => setSort(column.column_name)}
                    on:keypress={() => setSort(column.column_name)}
                    tabindex="0"
                    role="button"
                  >
                    <div>
                      {column.column_name}
                    </div>
                    <div
                      class="table-cell-icon
										{sortType !== '' && sortColumn === column.column_name ? sortType : ''}"
                    >
                      {#if sortType === "asc" && sortColumn === column.column_name}
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          viewBox="0 0 320 512"
                        >
                          <!--! Font Awesome Free 6.4.0 by @fontawesome - https://fontawesome.com License - https://fontawesome.com/license (Commercial License) Copyright 2023 Fonticons, Inc. -->
                          <path
                            d="M182.6 470.6c-12.5 12.5-32.8 12.5-45.3 0l-128-128c-9.2-9.2-11.9-22.9-6.9-34.9s16.6-19.8 29.6-19.8H288c12.9 0 24.6 7.8 29.6 19.8s2.2 25.7-6.9 34.9l-128 128z"
                          />
                        </svg>
                      {:else if sortType === "desc" && sortColumn === column.column_name}
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          viewBox="0 0 320 512"
                        >
                          <!--! Font Awesome Free 6.4.0 by @fontawesome - https://fontawesome.com License - https://fontawesome.com/license (Commercial License) Copyright 2023 Fonticons, Inc. -->
                          <path
                            d="M182.6 41.4c-12.5-12.5-32.8-12.5-45.3 0l-128 128c-9.2 9.2-11.9 22.9-6.9 34.9s16.6 19.8 29.6 19.8H288c12.9 0 24.6-7.8 29.6-19.8s2.2-25.7-6.9-34.9l-128-128z"
                          />
                        </svg>
                      {/if}
                    </div>
                  </div>
                </th>
              {/if}
            {/each}
          </tr>
          {#await recordsPromise then records}
            {#each records as record}
              <tr class="table-row">
                {#each columns as column, index}
                  {#if !activeColumns[index].is_hide}
                    <td class="table-cell">
                      {#if record[column.column_name] == null}
                        <span style="color: #999;">NULL</span>
                      {:else}
                        {record[column.column_name]}
                      {/if}
                    </td>
                  {/if}
                {/each}
              </tr>
            {/each}
          {/await}
        </table>
      </div>
    </div>
  {/await}
</div>

<style>
  .configuration {
    width: 100vw;
    display: flex;
    justify-content: space-between;
    padding: 0.75rem 1rem;
    font-size: 0.9rem;
  }

  .filter-section {
    display: flex;
    align-items: center;
    gap: 0.7rem;
  }

  .filter-btn,
  .column-btn {
    padding: 0.3rem 0.5rem;
    display: flex;
    gap: 0.3rem;
    align-items: center;
  }

  .column-toggle-wrapper {
    position: relative;
  }

  .column-toggle {
    display: none;
    position: absolute;
    top: 110%;
    left: -2%;
    width: 250px;
    padding: 0.5rem;
    border: 1px solid #777;
    border-radius: 0.3rem;
    background-color: #d9d9d9;
    max-height: 300px;
    overflow: scroll;
  }

  .column-toggle.is-active {
    display: block;
  }

  .column-toggle-list {
    display: flex;
    flex-direction: column;
  }

  .column-toggle-item {
    display: flex;
    align-items: center;
    text-align: left;
    padding: 0.25rem 0.5rem;
    gap: 0.3rem;
  }

  .column-toggle-icon {
    display: flex;
    align-items: center;
  }

  .filter-setting {
    font-size: 0.9rem;
    padding: 0 1rem 1rem 1rem;
  }

  .filter-setting-item {
    display: flex;
    justify-content: space-between;
    gap: 0.5rem;
    padding-bottom: 0.5rem;
  }

  .filter-setting-item-left {
    width: 100%;
    display: flex;
    gap: 0.5rem;
  }

  .filter-item-select-wrapper {
    position: relative;
    display: block;
  }

  .filter-item-field,
  .filter-item-operator {
    /* hide the default arrow icon */
    -webkit-appearance: none;
    -moz-appearance: none;
    appearance: none;

    padding: 0.3rem;
    border: 1px solid #777;
    border-radius: 0.3rem;
    cursor: pointer;
  }

  .filter-item-select-icon {
    position: absolute;
    top: 50%;
    right: 10px;
    transform: translateY(-50%);
    pointer-events: none; /* ensure the filter icon doesn't interfere with the select functionality */
    height: 0.9rem;
  }

  .filter-item-value {
    width: 100%;
    border: 1px solid #777;
    border-radius: 0.3rem;
    padding: 0 0.4rem;
    display: block;
  }

  .filter-item-btn {
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .filter-setting-submit {
    display: flex;
    gap: 0.5rem;
    justify-content: end;
  }

  .filter-setting-submit-btn {
    padding: 0.2rem 0.4rem;
  }

  .pagination-section,
  .pagination {
    display: flex;
    column-gap: 0.5rem;
    align-items: center;
  }

  .pagination-input input {
    width: 5.5rem;
    padding: 0 0.5rem 0 0.5rem;
    height: 30px;
    border-radius: 5px;
    border: 1px solid #888;
  }

  /* hide default up/down button in number input  */
  /* Chrome, Safari, Edge, Opera */
  .pagination-input input::-webkit-outer-spin-button,
  .pagination-input input::-webkit-inner-spin-button {
    -webkit-appearance: none;
    margin: 0;
  }
  /* Firefox */
  .pagination-input input[type="number"] {
    appearance: textfield;
  }

  .pagination-prev-btn,
  .pagination-next-btn {
    width: 30px;
    height: 30px;
    padding: 0;
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .pagination-prev-btn svg,
  .pagination-next-btn svg {
    height: 1.1rem;
  }

  .table-container {
    width: 100vw;
    overflow: auto;
  }

  .table-container > div {
    min-width: max-content;
    padding: 0 1rem 2rem 1rem;
  }

  .table {
    min-width: max-content;
    font-size: 0.9rem;
  }

  .table-row:hover {
    background-color: #d9d9d9;
  }

  .table-cell {
    font-size: 13.5px;
    text-align: left;
    padding: 3px 9px;
    border: 1px #777 solid;
    cursor: default;
    user-select: none;
  }

  .table-header .table-cell {
    padding: 0;
    font-size: 14px;
  }

  .table-cell-btn {
    display: flex;
    align-items: center;
    padding: 3px 9px;
    column-gap: 0.3rem;
    cursor: pointer;
    user-select: none;
  }

  .table-cell-icon {
    position: relative;
    width: 12px;
    height: 16px;
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .table-cell-icon svg {
    position: absolute;
    width: 10px;
  }

  .table-cell-icon.asc svg {
    bottom: 4px;
  }

  .table-cell-icon.desc svg {
    top: 4px;
  }

  svg {
    fill: #888;
  }
</style>
