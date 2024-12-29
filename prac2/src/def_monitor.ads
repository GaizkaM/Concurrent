-- def_monitor.ads
-----
-- Monitor del tipus Productor-Consumidor adaptat per a la resolució del exercici
-----
package def_monitor is

   protected type Maquina is
      -- Bloqueigs i desbloqueigs del monitor
      entry clientLock;
      procedure clientUnlock;
      entry reposadorLock;
      procedure reposadorUnlock;
   private
      clients : integer := 0;
      reposant : boolean := false;
   end Maquina;

end def_monitor;
