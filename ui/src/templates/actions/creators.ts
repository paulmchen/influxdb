// Types
import {CommunityTemplate} from 'src/types'
import {InstalledStack} from 'src/types'

export const SET_COMMUNITY_TEMPLATE_TO_INSTALL =
  'SET_COMMUNITY_TEMPLATE_TO_INSTALL'
export const TOGGLE_TEMPLATE_RESOURCE_INSTALL =
  'TOGGLE_TEMPLATE_RESOURCE_INSTALL'

export const SET_STACKS = 'SET_STACKS'
export const DELETE_STACKS = 'DELETE_STACKS'

export type Action =
  | ReturnType<typeof setCommunityTemplateToInstall>
  | ReturnType<typeof toggleTemplateResourceInstall>
  | ReturnType<typeof setStacks>
  | ReturnType<typeof removeStack>

// Action Creators
export const setCommunityTemplateToInstall = (template: CommunityTemplate) =>
  ({
    type: SET_COMMUNITY_TEMPLATE_TO_INSTALL,
    template,
  } as const)

export const toggleTemplateResourceInstall = (
  resourceType: string,
  templateMetaName: string,
  shouldInstall: boolean
) =>
  ({
    type: TOGGLE_TEMPLATE_RESOURCE_INSTALL,
    resourceType,
    templateMetaName,
    shouldInstall,
  } as const)

export const setStacks = (stacks: InstalledStack[]) =>
  ({
    type: SET_STACKS,
    stacks,
  } as const)

export const removeStack = (stackID: string) =>
  ({
    type: DELETE_STACKS,
    stackID,
  } as const)
