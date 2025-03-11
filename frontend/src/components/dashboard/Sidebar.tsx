'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { useAuth } from '@/components/auth/AuthProvider';

export default function Sidebar() {
  const pathname = usePathname();
  const { user } = useAuth();

  // Define navigation items
  const navigation = [
    { name: 'Dashboard', href: '/dashboard', icon: '📊' },
    { name: 'Contacts', href: '/contacts', icon: '👥' },
    { name: 'Bible Studies', href: '/studies', icon: '📚' },
    { name: 'Rooms', href: '/rooms', icon: '🏢' },
  ];
  
  // Admin-only navigation items
  const adminNavigation = [
    { name: 'Users', href: '/users', icon: '👤' },
    { name: 'Settings', href: '/settings', icon: '⚙️' },
  ];

  return (
    <div className="w-64 bg-gray-800 text-white">
      <div className="p-6">
        <h2 className="text-xl font-bold">Fruit Management</h2>
      </div>
      
      <nav className="mt-5">
        <div className="px-4 py-2 text-xs text-gray-400 uppercase">
          Main
        </div>
        
        <ul>
          {navigation.map((item) => (
            <li key={item.name}>
              <Link
                href={item.href}
                className={\`flex items-center px-6 py-3 \${
                  pathname?.startsWith(item.href)
                    ? 'bg-gray-900 text-white'
                    : 'text-gray-300 hover:bg-gray-700'
                }\`}
              >
                <span className="mr-3">{item.icon}</span>
                {item.name}
              </Link>
            </li>
          ))}
        </ul>
        
        {user?.role === 'admin' && (
          <>
            <div className="px-4 py-2 mt-5 text-xs text-gray-400 uppercase">
              Administration
            </div>
            <ul>
              {adminNavigation.map((item) => (
                <li key={item.name}>
                  <Link
                    href={item.href}
                    className={\`flex items-center px-6 py-3 \${
                      pathname?.startsWith(item.href)
                        ? 'bg-gray-900 text-white'
                        : 'text-gray-300 hover:bg-gray-700'
                    }\`}
                  >
                    <span className="mr-3">{item.icon}</span>
                    {item.name}
                  </Link>
                </li>
              ))}
            </ul>
          </>
        )}
      </nav>
    </div>
  );
}