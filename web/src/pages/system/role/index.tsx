import React, { useEffect, useRef, useState } from "react";
import { DeleteOutlined, EditOutlined, PlusOutlined } from "@ant-design/icons";
import type { ActionType, ProColumns } from "@ant-design/pro-components";
import {
  ModalForm,
  ProFormDigit,
  ProFormRadio,
  ProFormText,
  ProFormTreeSelect,
  ProTable,
} from "@ant-design/pro-components";
import { Button, Form, Popconfirm, Tag, TreeSelect, message } from "antd";
import { getMenuList, type MenuItem } from "@/api/menu";
import {
  createRole,
  deleteRole,
  getRoleDetail,
  getRoleList,
  updateRole,
  type RoleItem,
  type RolePayload,
} from "@/api/role";

const RoleManagement: React.FC = () => {
  const actionRef = useRef<ActionType>();
  const [modalVisible, setModalVisible] = useState(false);
  const [currentRow, setCurrentRow] = useState<RoleItem | null>(null);
  const [menuTree, setMenuTree] = useState<MenuItem[]>([]);
  const [form] = Form.useForm<RolePayload>();

  const fetchMenus = async () => {
    try {
      const data = await getMenuList();
      setMenuTree(data as MenuItem[]);
    } catch (error) {
      message.error("获取菜单失败");
    }
  };

  const fetchData = async () => {
    try {
      const data = await getRoleList();
      return {
        data: data as RoleItem[],
        success: true,
      };
    } catch (error) {
      return {
        data: [],
        success: false,
      };
    }
  };

  const openCreate = () => {
    setCurrentRow(null);
    form.resetFields();
    setModalVisible(true);
  };

  const openEdit = async (record: RoleItem) => {
    setCurrentRow(record);
    setModalVisible(true);
    try {
      const detail = await getRoleDetail(record.id);
      form.setFieldsValue({
        name: detail.name,
        key: detail.key,
        sort: detail.sort,
        status: detail.status,
        menuIds: detail.menuIds || [],
      });
    } catch (error) {
      message.error("加载角色详情失败");
    }
  };

  const handleDelete = async (id: number) => {
    try {
      await deleteRole(id);
      message.success("删除成功");
      actionRef.current?.reload();
    } catch (error) {
      message.error("删除失败");
    }
  };

  useEffect(() => {
    fetchMenus();
  }, []);

  const columns: ProColumns<RoleItem>[] = [
    {
      title: "角色名称",
      dataIndex: "name",
      width: 200,
    },
    {
      title: "角色标识",
      dataIndex: "key",
      width: 160,
      render: (_, record) => <Tag>{record.key}</Tag>,
    },
    {
      title: "排序",
      dataIndex: "sort",
      width: 80,
    },
    {
      title: "状态",
      dataIndex: "status",
      width: 100,
      render: (_, record) =>
        record.status === 1 ? <Tag color="green">启用</Tag> : <Tag>停用</Tag>,
    },
    {
      title: "操作",
      valueType: "option",
      width: 160,
      render: (_, record) => [
        <Button
          key="edit"
          type="link"
          icon={<EditOutlined />}
          onClick={() => openEdit(record)}
        >
          编辑
        </Button>,
        <Popconfirm
          key="delete"
          title="确定删除吗？"
          onConfirm={() => handleDelete(record.id)}
        >
          <Button type="link" danger icon={<DeleteOutlined />}>
            删除
          </Button>
        </Popconfirm>,
      ],
    },
  ];

  return (
    <>
      <ProTable<RoleItem>
        headerTitle="角色列表"
        actionRef={actionRef}
        rowKey="id"
        search={false}
        toolBarRender={() => [
          <Button
            type="primary"
            key="primary"
            icon={<PlusOutlined />}
            onClick={openCreate}
          >
            新建
          </Button>,
        ]}
        request={fetchData}
        columns={columns}
        pagination={false}
      />

      <ModalForm
        title={currentRow?.id ? "编辑角色" : "新建角色"}
        width="600px"
        form={form}
        visible={modalVisible}
        onVisibleChange={setModalVisible}
        onFinish={async (value) => {
          const payload: RolePayload = {
            name: value.name,
            key: value.key,
            sort: value.sort || 0,
            status: value.status ?? 1,
            menuIds: value.menuIds || [],
          };
          try {
            if (currentRow?.id) {
              await updateRole(currentRow.id, payload);
              message.success("更新成功");
            } else {
              await createRole(payload);
              message.success("创建成功");
            }
            actionRef.current?.reload();
            return true;
          } catch (error) {
            message.error("保存失败");
            return false;
          }
        }}
        modalProps={{
          destroyOnClose: true,
        }}
      >
        <ProFormText
          name="name"
          label="角色名称"
          placeholder="例如：运营"
          rules={[{ required: true }]}
        />
        <ProFormText
          name="key"
          label="角色标识"
          placeholder="例如：operator"
          rules={[{ required: true }]}
        />
        <ProFormDigit width="md" name="sort" label="排序" min={0} />
        <ProFormRadio.Group
          name="status"
          label="状态"
          initialValue={1}
          options={[
            { label: "启用", value: 1 },
            { label: "停用", value: 0 },
          ]}
        />
        <ProFormTreeSelect
          name="menuIds"
          label="菜单权限"
          placeholder="请选择菜单权限"
          fieldProps={{
            treeData: menuTree,
            fieldNames: { label: "title", value: "id", children: "children" },
            treeCheckable: true,
            showCheckedStrategy: TreeSelect.SHOW_PARENT,
            allowClear: true,
          }}
        />
      </ModalForm>
    </>
  );
};

export default RoleManagement;
