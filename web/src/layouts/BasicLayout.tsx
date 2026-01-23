import { ProLayout } from '@ant-design/pro-components';
import { Outlet, Link } from 'react-router-dom';

const BasicLayout = () => {
  return (
    <ProLayout
      route={{
        path: '/',
        routes: [
          {
            path: '/dashboard',
            name: 'Dashboard',
            icon: 'crown',
          },
        ],
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
