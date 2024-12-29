-- def_monitor.adb
-----
-- Cos del Monitor amb els bloqueigs i desbloqueigs
-----
package body def_monitor is

   protected body Maquina is
      entry clientLock when not reposant is
      begin
         clients := clients + 1;
      end clientLock;

      procedure clientUnlock is
      begin
         clients := clients - 1;
      end clientUnlock;

      entry reposadorLock when (clients = 0) and (not reposant) is
      begin
         reposant := true;
      end reposadorLock;

      procedure reposadorUnlock is
      begin
         reposant := false;
      end reposadorUnlock;

   end Maquina;

end def_monitor;
