import { useEffect, useState, useMemo } from 'react';
import { cn } from '@/lib/utils';
import {
  User2,
  Building2,
  Mail,
  Key,
  Plus,
  Trash2,
  Loader2,
  ChevronDown,
  ChevronUp,
} from 'lucide-react';
import { userService } from '@/services/users';
import type { User } from '@/types/user';
import { Skeleton } from '@/components/ui/skeleton';
import { Button } from '@/components/ui/button';
import { UserDialog } from '@/components/users/user-dialog';
import { InputSearch } from '@/components/ui/input-search';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog';

type SortConfig = { key: keyof User | null; direction: 'asc' | 'desc' };

const UserTableSkeleton = () => (
  <div className="space-y-4">
    {/* Skeleton for the header section */}
    <div className="flex justify-between items-center">
      <Skeleton className="h-10 w-32" />
      <Skeleton className="h-10 w-32" />
    </div>

    {/* Skeleton for the search and filter section */}
    <div className="flex flex-col sm:flex-row gap-4 mb-4">
      <Skeleton className="h-10 flex-1" />
      <Skeleton className="h-10 w-48" />
    </div>

    {/* Skeleton for the table */}
    <div className="rounded-md border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>
              <div className="flex items-center gap-2">
                <User2 className="h-4 w-4 text-muted-foreground" />
                <span>Name</span>
              </div>
            </TableHead>
            <TableHead>
              <div className="flex items-center gap-2">
                <Building2 className="h-4 w-4 text-muted-foreground" />
                <span>Company</span>
              </div>
            </TableHead>
            <TableHead>
              <div className="flex items-center gap-2">
                <Mail className="h-4 w-4 text-muted-foreground" />
                <span>Email</span>
              </div>
            </TableHead>
            <TableHead>
              <div className="flex items-center gap-2">
                <Key className="h-4 w-4 text-muted-foreground" />
                <span>API Key</span>
              </div>
            </TableHead>
            <TableHead>Status</TableHead>
            <TableHead>Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {[...Array(5)].map((_, index) => (
            <TableRow key={index}>
              <TableCell>
                <Skeleton className="h-5 w-[120px]" />
              </TableCell>
              <TableCell>
                <Skeleton className="h-5 w-[100px]" />
              </TableCell>
              <TableCell>
                <Skeleton className="h-5 w-[180px]" />
              </TableCell>
              <TableCell>
                <Skeleton className="h-5 w-[140px]" />
              </TableCell>
              <TableCell>
                <Skeleton className="h-6 w-16 rounded-full" />
              </TableCell>
              <TableCell>
                <div className="flex items-center gap-2">
                  <Skeleton className="h-8 w-16" />
                  <Skeleton className="h-8 w-8" />
                </div>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  </div>
);

export function Users() {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [dialogOpen, setDialogOpen] = useState(false);
  const [editingUser, setEditingUser] = useState<User | undefined>(undefined);
  const [deletingUser, setDeletingUser] = useState<User | undefined>(undefined);
  const [actionLoading, setActionLoading] = useState({
    create: false,
    edit: false,
    delete: false,
  });
  const [filters, setFilters] = useState({
    search: '',
    status: 'all' as 'all' | 'active' | 'inactive',
  });
  const [sortConfig, setSortConfig] = useState<SortConfig>({
    key: null,
    direction: 'asc',
  });

  const fetchUsers = async () => {
    try {
      setLoading(true);
      const response = await userService.list();
      setUsers(response.users);
      setError(null);
    } catch (err) {
      setError('Failed to fetch users');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchUsers();
  }, []);

  const handleCreateUser = async (data: Omit<User, 'id' | 'apiKey'>) => {
    try {
      setActionLoading((prev) => ({ ...prev, create: true }));
      const newUser = await userService.create(data);
      setUsers((prev) => [...prev, newUser]);
      setDialogOpen(false);
    } catch (err) {
      console.error('Error creating user:', err);
    } finally {
      setActionLoading((prev) => ({ ...prev, create: false }));
    }
  };

  const handleEditUser = async (data: Omit<User, 'id' | 'apiKey'>) => {
    if (!editingUser) return;
    try {
      setActionLoading((prev) => ({ ...prev, edit: true }));
      const updatedUser = await userService.update(editingUser.id, data);
      setUsers((prev) => prev.map((user) => (user.id === editingUser.id ? updatedUser : user)));
      setDialogOpen(false);
    } catch (err) {
      console.error('Error updating user:', err);
    } finally {
      setActionLoading((prev) => ({ ...prev, edit: false }));
    }
  };

  const handleDeleteUser = async () => {
    if (!deletingUser) return;
    try {
      setActionLoading((prev) => ({ ...prev, delete: true }));
      await userService.delete(deletingUser.id);
      setUsers(users.filter((user) => user.id !== deletingUser.id));
      setDeletingUser(undefined);
    } catch (err) {
      console.error('Error deleting user:', err);
    } finally {
      setActionLoading((prev) => ({ ...prev, delete: false }));
    }
  };

  const filteredUsers = users.filter((user) => {
    const matchesSearch =
      filters.search === '' ||
      user.name.toLowerCase().includes(filters.search.toLowerCase()) ||
      user.company.toLowerCase().includes(filters.search.toLowerCase()) ||
      user.email.toLowerCase().includes(filters.search.toLowerCase());
    const matchesStatus = filters.status === 'all' || user.status === filters.status;
    return matchesSearch && matchesStatus;
  });

  const sortedUsers = useMemo(() => {
    const sorted = [...filteredUsers];
    if (sortConfig.key) {
      sorted.sort((a, b) => {
        if (a[sortConfig.key!] < b[sortConfig.key!]) return sortConfig.direction === 'asc' ? -1 : 1;
        if (a[sortConfig.key!] > b[sortConfig.key!]) return sortConfig.direction === 'asc' ? 1 : -1;
        return 0;
      });
    }
    return sorted;
  }, [filteredUsers, sortConfig]);

  const handleSort = (key: keyof User) => {
    setSortConfig({
      key,
      direction: sortConfig.key === key && sortConfig.direction === 'asc' ? 'desc' : 'asc',
    });
  };

  const SortIndicator = ({ column }: { column: keyof User }) => {
    if (sortConfig.key !== column) return null;
    return sortConfig.direction === 'asc' ? (
      <ChevronUp className="h-4 w-4" />
    ) : (
      <ChevronDown className="h-4 w-4" />
    );
  };

  if (loading) return <UserTableSkeleton />;
  if (error)
    return (
      <div className="flex items-center justify-center text-destructive">
        <span>{error}</span>
      </div>
    );

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <h2 className="text-3xl font-bold tracking-tight">Users</h2>
        <Button onClick={() => setDialogOpen(true)} disabled={actionLoading.create}>
          {actionLoading.create ? (
            <Loader2 className="mr-2 h-4 w-4 animate-spin" />
          ) : (
            <Plus className="mr-2 h-4 w-4" />
          )}
          Add User
        </Button>
      </div>

      <div className="flex flex-col sm:flex-row gap-4 mb-4">
        <div className="flex-1">
          <InputSearch
            value={filters.search}
            onChange={(value) => setFilters((prev) => ({ ...prev, search: value }))}
            placeholder="Search by name, company or email..."
          />
        </div>
        <div className="w-full sm:w-48">
          <select
            className="w-full h-10 px-3 rounded-md border border-input bg-background text-sm"
            value={filters.status}
            onChange={(e) =>
              setFilters((prev) => ({
                ...prev,
                status: e.target.value as 'all' | 'active' | 'inactive',
              }))
            }
          >
            <option value="all">All Status</option>
            <option value="active">Active</option>
            <option value="inactive">Inactive</option>
          </select>
        </div>
      </div>

      <div className="rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="cursor-pointer" onClick={() => handleSort('name')}>
                <div className="flex items-center gap-2">
                  <User2 className="h-4 w-4" />
                  Name
                  <SortIndicator column="name" />
                </div>
              </TableHead>
              <TableHead className="cursor-pointer" onClick={() => handleSort('company')}>
                <div className="flex items-center gap-2">
                  <Building2 className="h-4 w-4" />
                  Company
                  <SortIndicator column="company" />
                </div>
              </TableHead>
              <TableHead className="cursor-pointer" onClick={() => handleSort('email')}>
                <div className="flex items-center gap-2">
                  <Mail className="h-4 w-4" />
                  Email
                  <SortIndicator column="email" />
                </div>
              </TableHead>
              <TableHead>
                <div className="flex items-center gap-2">
                  <Key className="h-4 w-4" />
                  API Key
                </div>
              </TableHead>
              <TableHead className="cursor-pointer" onClick={() => handleSort('status')}>
                <div className="flex items-center gap-2">
                  Status
                  <SortIndicator column="status" />
                </div>
              </TableHead>
              <TableHead>Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {sortedUsers.map((user) => (
              <TableRow key={user.id}>
                <TableCell className="font-medium">{user.name}</TableCell>
                <TableCell>{user.company}</TableCell>
                <TableCell>{user.email}</TableCell>
                <TableCell className="font-mono">{user.apiKey}</TableCell>
                <TableCell>
                  <span
                    className={cn(
                      'inline-flex items-center rounded-full px-2 py-1 text-xs font-medium',
                      user.status === 'active'
                        ? 'bg-green-50 text-green-700 border border-green-200'
                        : 'bg-red-50 text-red-700 border border-red-200'
                    )}
                  >
                    {user.status}
                  </span>
                </TableCell>
                <TableCell>
                  <div className="flex items-center gap-2">
                    <Button
                      variant="ghost"
                      onClick={() => {
                        setEditingUser(user);
                        setDialogOpen(true);
                      }}
                      disabled={actionLoading.edit && editingUser?.id === user.id}
                    >
                      {actionLoading.edit && editingUser?.id === user.id ? (
                        <Loader2 className="h-4 w-4 animate-spin" />
                      ) : (
                        'Edit'
                      )}
                    </Button>
                    <Button
                      variant="ghost"
                      onClick={() => setDeletingUser(user)}
                      disabled={actionLoading.delete && deletingUser?.id === user.id}
                      className="text-red-600 hover:text-red-800"
                    >
                      {actionLoading.delete && deletingUser?.id === user.id ? (
                        <Loader2 className="h-4 w-4 animate-spin" />
                      ) : (
                        <Trash2 className="h-4 w-4" />
                      )}
                    </Button>
                  </div>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>

      <UserDialog
        open={dialogOpen}
        onOpenChange={(open) => {
          setDialogOpen(open);
          if (!open) setEditingUser(undefined);
        }}
        onSubmit={editingUser ? handleEditUser : handleCreateUser}
        user={editingUser}
        loading={editingUser ? actionLoading.edit : actionLoading.create}
      />

      <AlertDialog open={!!deletingUser} onOpenChange={() => setDeletingUser(undefined)}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Are you sure?</AlertDialogTitle>
            <AlertDialogDescription>
              This action cannot be undone. This will permanently delete the user and remove their
              data from our servers.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction
              onClick={handleDeleteUser}
              className="bg-red-600 hover:bg-red-700"
              disabled={actionLoading.delete}
            >
              {actionLoading.delete ? <Loader2 className="h-4 w-4 animate-spin" /> : 'Delete'}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}
