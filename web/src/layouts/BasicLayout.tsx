import { ProLayout } from "@ant-design/pro-components";
import { Outlet, Link } from "react-router-dom";
import useMenuStore from "@/store/useMenuStore";
import { buildLayoutMenus } from "@/router/routeUtils";

const BasicLayout = () => {
  const menus = useMenuStore((state) => state.menus);
  const layoutMenus = buildLayoutMenus(menus);

  return (
    <ProLayout
      route={{
        path: "/",
        routes: layoutMenus,
      }}
      menuItemRender={(menuItemProps, defaultDom) => {
        if (menuItemProps.isUrl || !menuItemProps.path) {
          return defaultDom;
        }
        return <Link to={menuItemProps.path}>{defaultDom}</Link>;
      }}
    >
      <Outlet />
    </ProLayout>
  );
};

export default BasicLayout;
