<script lang="ts">
  import Table from '../../components/global/table/Table.svelte';
  import HeadCol from '../../components/global/table/HeadCol.svelte';
  import TableRow from '../../components/global/table/TableRow.svelte';
  import RowCol from '../../components/global/table/RowCol.svelte';

  import LL from '../../i18n/i18n-svelte';
  import SolidButton from '../../components/global/SolidButton.svelte';
  import UserAvatar from '../../components/user/UserAvatar.svelte';
  import CountryFlag from '../../components/user/CountryFlag.svelte';
  import HollowButton from '../../components/global/HollowButton.svelte';
  import AddUser from '../../components/team/AddUser.svelte';
  import UpdateUser from '../../components/team/UpdateUser.svelte';
  import DeleteConfirmation from '../../components/global/DeleteConfirmation.svelte';

  export let xfetch;
  export let notifications;
  export let eventTag;
  export let teamPrefix = '';
  export let isAdmin = false;
  export let pageType = '';
  export let users = [];
  export let getUsers = () => {};
  let showAddUser = false;
  let showUpdateUser = false;
  let updateUser = {};
  let showRemoveUser = false;
  let removeUserId = null;

  function handleUserAdd(email, role) {
    const body = {
      email,
      role,
    };

    xfetch(`${teamPrefix}/users`, { body })
      .then(function () {
        eventTag(`${pageType}_add_user`, 'engagement', 'success');
        toggleAddUser();
        notifications.success($LL.userAddSuccess());
        getUsers();
      })
      .catch(function () {
        notifications.danger($LL.userAddError());
        eventTag(`${pageType}_add_user`, 'engagement', 'failure');
      });
  }

  function handleUserUpdate(userId, role) {
    const body = {
      role,
    };

    xfetch(`${teamPrefix}/users/${userId}`, { body, method: 'PUT' })
      .then(function () {
        eventTag(`${pageType}_update_user`, 'engagement', 'success');
        toggleUpdateUser({})();
        notifications.success($LL.userUpdateSuccess());
        getUsers();
      })
      .catch(function () {
        notifications.danger($LL.userUpdateError());
        eventTag(`${pageType}_update_user`, 'engagement', 'failure');
      });
  }

  function handleUserRemove() {
    xfetch(`${teamPrefix}/users/${removeUserId}`, { method: 'DELETE' })
      .then(function () {
        eventTag(`${pageType}_remove_user`, 'engagement', 'success');
        toggleRemoveUser(null)();
        notifications.success($LL.userRemoveSuccess());
        getUsers();
      })
      .catch(function () {
        notifications.danger($LL.userRemoveError());
        eventTag(`${pageType}_remove_user`, 'engagement', 'failure');
      });
  }

  function toggleAddUser() {
    showAddUser = !showAddUser;
  }

  const toggleUpdateUser = user => () => {
    updateUser = user;
    showUpdateUser = !showUpdateUser;
  };

  const toggleRemoveUser = userId => () => {
    showRemoveUser = !showRemoveUser;
    removeUserId = userId;
  };
</script>

<div class="w-full">
  <div class="flex w-full">
    <div class="w-4/5">
      <h2
        class="text-2xl font-semibold font-rajdhani uppercase mb-4 dark:text-white"
      >
        {$LL.users()}
      </h2>
    </div>
    <div class="w-1/5">
      <div class="text-right">
        {#if isAdmin}
          <SolidButton onClick="{toggleAddUser}" testid="user-add">
            {$LL.userAdd()}
          </SolidButton>
        {/if}
      </div>
    </div>
  </div>

  <Table>
    <tr slot="header">
      <HeadCol>
        {$LL.name()}
      </HeadCol>
      <HeadCol>
        {$LL.email()}
      </HeadCol>
      <HeadCol>
        {$LL.role()}
      </HeadCol>
      <HeadCol type="action">
        <span class="sr-only">{$LL.actions()}</span>
      </HeadCol>
    </tr>
    <tbody slot="body" let:class="{className}" class="{className}">
      {#each users as user, i}
        <TableRow itemIndex="{i}">
          <RowCol>
            <div class="flex items-center">
              <div class="flex-shrink-0 h-10 w-10">
                <UserAvatar
                  warriorId="{user.id}"
                  avatar="{user.avatar}"
                  gravatarHash="{user.gravatarHash}"
                  userName="{user.name}"
                  width="48"
                  class="h-10 w-10 rounded-full"
                />
              </div>
              <div class="ms-4">
                <div class="font-medium text-gray-900 dark:text-gray-200">
                  <span data-testid="user-name">{user.name}</span>
                  {#if user.country}
                    &nbsp;
                    <CountryFlag
                      country="{user.country}"
                      additionalClass="inline-block"
                      width="32"
                      height="24"
                    />
                  {/if}
                </div>
              </div>
            </div>
          </RowCol>
          <RowCol>
            <span data-testid="user-email">{user.email}</span>
          </RowCol>
          <RowCol>
            <div class="text-sm text-gray-500 dark:text-gray-300">
              {user.role}
            </div>
          </RowCol>
          <RowCol type="action">
            {#if isAdmin}
              <HollowButton onClick="{toggleUpdateUser(user)}" color="blue">
                {$LL.edit()}
              </HollowButton>
              <HollowButton onClick="{toggleRemoveUser(user.id)}" color="red">
                {$LL.remove()}
              </HollowButton>
            {/if}
          </RowCol>
        </TableRow>
      {/each}
    </tbody>
  </Table>
</div>

{#if showAddUser}
  <AddUser toggleAdd="{toggleAddUser}" handleAdd="{handleUserAdd}" />
{/if}

{#if showUpdateUser}
  <UpdateUser
    toggleUpdate="{toggleUpdateUser({})}"
    handleUpdate="{handleUserUpdate}"
    userId="{updateUser.id}"
    userEmail="{updateUser.email}"
    role="{updateUser.role}"
  />
{/if}

{#if showRemoveUser}
  <DeleteConfirmation
    toggleDelete="{toggleRemoveUser(null)}"
    handleDelete="{handleUserRemove}"
    permanent="{false}"
    confirmText="{$LL.removeUserConfirmText()}"
    confirmBtnText="{$LL.removeUser()}"
  />
{/if}
