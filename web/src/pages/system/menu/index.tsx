import React, { useRef, useState } from "react"
import { DeleteOutlined, EditOutlined, PlusOutlined } from "@ant-design/icons"
import type { ActionType, ProColumns } from "@ant-design/pro-components"
import { ProTable } from "@ant-design/pro-components"
import {
  getMenuList,
  createMenu,
  updateMenu,
  deleteMenu,
  type MenuItem,
} from "@/api/menu"
import { message, Button, Tag, Popconfirm } from "antd"
import MenuModal from "./components/menuModal" // 引入拆分后的组件

const MenuManagement: React.FC = () => {
  const actionRef = useRef<ActionType | null>(null)
  const [modalVisible, setModalVisible] = useState(false)
  const [currentRow, setCurrentRow] = useState<Partial<MenuItem> | null>(null)
  const [menuTree, setMenuTree] = useState<MenuItem[]>([])

  const fetchData = async () => {
    try {
      const res = await getMenuList()
      const data = res as unknown as MenuItem[]
      setMenuTree(data)
      return {
        data,
        success: true,
      }
    } catch (error) {
      return {
        data: [],
        success: false,
      }
    }
  }

  const handleAdd = async (fields: MenuItem) => {
    try {
      await createMenu(fields)
      message.success("菜单创建成功")
      actionRef.current?.reload()
      setModalVisible(false)
    } catch (error) {
      message.error("添加失败")
    }
  }

  const handleUpdate = async (fields: MenuItem) => {
    if (!currentRow?.id) return
    try {
      await updateMenu(currentRow.id, fields)
      message.success("更新成功")
      setModalVisible(false)
      setCurrentRow(null)
      actionRef.current?.reload()
    } catch (error) {
      message.error("更新失败")
    }
  }

  const handleDelete = async (id: number) => {
    try {
      await deleteMenu(id)
      message.success("删除成功")
      actionRef.current?.reload()
    } catch (error) {
      message.error("删除失败")
    }
  }

  // 统一处理提交
  const handleSubmit = async (values: MenuItem) => {
    if (currentRow?.id) {
      await handleUpdate(values)
    } else {
      await handleAdd(values)
    }
  }

  const columns: ProColumns<MenuItem>[] = [
    {
      title: "菜单名称",
      dataIndex: "title",
      width: 200,
      fixed: "left",
    },
    {
      title: "图标",
      dataIndex: "icon",
      width: 80,
      render: (_, record) => (record.icon ? <Tag>{record.icon}</Tag> : "-"),
    },
    {
      title: "排序",
      dataIndex: "sort",
      width: 60,
    },
    {
      title: "类型",
      dataIndex: "type",
      valueEnum: {
        1: { text: "目录", status: "Processing" },
        2: { text: "菜单", status: "Success" },
        3: { text: "按钮", status: "Warning" },
      },
      width: 80,
    },
    {
      title: "路由路径",
      dataIndex: "path",
      ellipsis: true,
    },
    {
      title: "组件路径",
      dataIndex: "component",
      ellipsis: true,
    },
    {
      title: "权限标识",
      dataIndex: "perms",
      ellipsis: true,
    },
    {
      title: "操作",
      valueType: "option",
      width: 180,
      fixed: "right",
      render: (_, record) => [
        <Button
          key="edit"
          type="link"
          icon={<EditOutlined />}
          onClick={() => {
            setCurrentRow(record)
            setModalVisible(true)
          }}
        >
          编辑
        </Button>,
        <Button
          key="addSub"
          type="link"
          size="small"
          icon={<PlusOutlined />}
          onClick={() => {
            setCurrentRow({ pid: record.id }) // 设置父ID
            setModalVisible(true)
          }}
        >
          新增
        </Button>,
        <Popconfirm
          key="delete"
          title="确定删除吗？"
          description="删除后无法恢复，且子菜单可能也会受到影响。"
          onConfirm={() => handleDelete(record.id)}
        >
          <Button type="link" danger icon={<DeleteOutlined />}>
            删除
          </Button>
        </Popconfirm>,
      ],
    },
  ]

  return (
    <>
      <ProTable<MenuItem>
        headerTitle="菜单列表"
        actionRef={actionRef}
        rowKey="id"
        search={false}
        toolBarRender={() => [
          <Button
            type="primary"
            key="primary"
            icon={<PlusOutlined />}
            onClick={() => {
              setCurrentRow(null)
              setModalVisible(true)
            }}
          >
            新建
          </Button>,
        ]}
        request={fetchData}
        columns={columns}
        pagination={false}
      />

      <MenuModal
        visible={modalVisible}
        onVisibleChange={setModalVisible}
        initialValues={currentRow || {}}
        onFinish={async (value) => {
          if (currentRow?.id) {
            await handleUpdate(value as MenuItem)
          } else {
            await hanldeAdd(value as MenuItem)
          }
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
        />

        <ProFormTreeSelect
          name="pid"
          label="上级菜单"
          placeholder="请选择上级菜单（留空则为顶级）"
          fieldProps={{
            treeData: menuTree,
            fieldNames: { label: "title", value: "id", children: "children" },
            treeDefaultExpandAll: true,
          }}
        />

        <div style={{ display: "flex", gap: 16 }}>
          <ProFormText
            width="md"
            name="title"
            label="显示标题"
            placeholder="例如：系统管理"
            rules={[{ required: true }]}
          />
          <ProFormText
            width="md"
            name="name"
            label="路由名称"
            placeholder="例如：System"
          />
        </div>

        <div style={{ display: "flex", gap: 16 }}>
          <ProFormText
            width="md"
            name="icon"
            label="图标"
            placeholder="Antd 图标名称"
          />
          <ProFormDigit width="md" name="sort" label="排序" min={0} />
        </div>

        {/* 仅目录和菜单显示 */}
        <ProFormText name="path" label="路由路径" placeholder="例如：/system" />
        <ProFormText
          name="redirect"
          label="重定向路径"
          placeholder="例如：/system/menu"
        />

        {/* 仅菜单显示 */}
        <ProFormText
          name="component"
          label="组件路径"
          placeholder="例如：/system/menu/index"
        />

        <ProFormText
          name="perms"
          label="权限标识"
          placeholder="例如：sys:menu:add"
        />

        <div style={{ display: "flex", gap: 16 }}>
          <ProFormSwitch name="hidden" label="隐藏菜单" />
          <ProFormSwitch
            name="keepAlive"
            label="页面缓存"
            initialValue={true}
          />
        </div>
      </ModalForm>
    </>
  )
}

export default MenuManagement
