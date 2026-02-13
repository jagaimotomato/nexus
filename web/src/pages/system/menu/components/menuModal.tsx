import React from "react"
import {
  ModalForm,
  ProFormText,
  ProFormDigit,
  ProFormRadio,
  ProFormTreeSelect,
  ProFormSwitch,
  ProFormDependency,
} from "@ant-design/pro-components"
import type { MenuItem } from "@/api/menu"

export interface MenuModalProps {
  visible: boolean
  onVisibleChange: (visible: boolean) => void
  currentRow: Partial<MenuItem> | null
  menuTree: MenuItem[]
  onSubmit: (values: MenuItem) => Promise<void>
}

const MenuModal: React.FC<MenuModalProps> = (props) => {
  const { visible, onVisibleChange, currentRow, menuTree, onSubmit } = props

  return (
    <ModalForm
      title={currentRow?.id ? "编辑菜单" : "新建菜单"}
      width="600px"
      visible={visible}
      onVisibleChange={onVisibleChange}
      // 注意：这里使用了 currentRow，当它变化时 initialValues 不会自动重置，
      // 需要配合 modalProps.destroyOnClose 来确保每次打开都是新的
      initialValues={currentRow || { type: 1, sort: 0, keepAlive: true }}
      onFinish={async (value) => {
        await onSubmit(value as MenuItem)
        return true
      }}
      modalProps={{
        destroyOnClose: true,
      }}
    >
      <ProFormRadio.Group
        name="type"
        label="菜单类型"
        options={[
          { label: "目录", value: 1 },
          { label: "菜单", value: 2 },
          { label: "按钮", value: 3 },
        ]}
        rules={[{ required: true }]}
      />

      <ProFormTreeSelect
        name="pid"
        label="上级菜单"
        placeholder="请选择上级菜单（留空则为顶级）"
        fieldProps={{
          treeData: menuTree,
          fieldNames: { label: "title", value: "id", children: "children" },
          treeDefaultExpandAll: false,
        }}
      />

      <div style={{ display: "flex", gap: 16 }}>
        <ProFormText
          width="md"
          name="title"
          label="显示标题"
          placeholder="例如：系统管理"
          rules={[{ required: true, message: "请输入显示标题" }]}
        />
        <ProFormDigit
          width="md"
          name="sort"
          label="排序"
          min={0}
          tooltip="数值越小越靠前"
        />
      </div>

      {/* 字段联动逻辑 */}
      <ProFormDependency name={["type"]}>
        {({ type }) => {
          return (
            <>
              {/* 图标：仅目录(1)和菜单(2)需要 */}
              {(type === 1 || type === 2) && (
                <div style={{ display: "flex", gap: 16 }}>
                  <ProFormText
                    width="md"
                    name="icon"
                    label="图标"
                    placeholder="Antd 图标名称"
                  />
                  <ProFormText
                    width="md"
                    name="name"
                    label="路由名称"
                    placeholder="例如：System"
                    tooltip="用于 keep-alive 缓存，必须唯一"
                  />
                </div>
              )}

              {/* 路由路径：仅目录(1)和菜单(2)需要 */}
              {(type === 1 || type === 2) && (
                <div style={{ display: "flex", gap: 16 }}>
                  <ProFormText
                    width="md"
                    name="path"
                    label="路由路径"
                    placeholder="例如：/system"
                    rules={[{ required: true, message: "请输入路由路径" }]}
                  />
                  {/* 重定向：通常仅目录(1)需要 */}
                  {type === 1 && (
                    <ProFormText
                      width="md"
                      name="redirect"
                      label="重定向"
                      placeholder="例如：/system/user"
                    />
                  )}
                </div>
              )}

              {/* 组件路径：仅菜单(2)需要 */}
              {type === 2 && (
                <ProFormText
                  name="component"
                  label="组件路径"
                  placeholder="例如：/system/menu/index"
                  tooltip="src/pages 下的文件路径"
                  rules={[{ required: true, message: "请输入组件路径" }]}
                />
              )}

              {/* 权限标识：菜单(2)和按钮(3)需要 */}
              {(type === 2 || type === 3) && (
                <ProFormText
                  name="perms"
                  label="权限标识"
                  placeholder="例如：sys:menu:add"
                />
              )}

              {/* 开关选项：仅目录(1)和菜单(2)需要 */}
              {(type === 1 || type === 2) && (
                <div style={{ display: "flex", gap: 16, marginTop: 16 }}>
                  <ProFormSwitch
                    name="hidden"
                    label="隐藏菜单"
                    checkedChildren="是"
                    unCheckedChildren="否"
                  />
                  {type === 2 && (
                    <ProFormSwitch
                      name="keepAlive"
                      label="页面缓存"
                      checkedChildren="是"
                      unCheckedChildren="否"
                      initialValue={true}
                    />
                  )}
                </div>
              )}
            </>
          )
        }}
      </ProFormDependency>
    </ModalForm>
  )
}

export default MenuModal
